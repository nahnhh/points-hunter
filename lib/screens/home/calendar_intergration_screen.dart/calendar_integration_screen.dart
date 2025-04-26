// courtesy of Kazutaka Yoshinaga
import 'package:clocky_app/api/api_client_provider.dart';
import 'package:clocky_app/api/user_notifier.dart';
import 'package:clocky_app/firebase_options.dart';
import 'package:clocky_app/services/google_calendar_auth_service.dart';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:google_sign_in/google_sign_in.dart';

class CalendarIntegrationScreen extends ConsumerWidget {
  const CalendarIntegrationScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('カレンダー連携'),
      ),
      body: Center(
        child: ElevatedButton(
          onPressed: () async {
            await _fetchTodayEvents(context, ref);
          },
          child: const Text('Googleカレンダーと連携する'),
        ),
      ),
    );
  }

  Future _fetchTodayEvents(BuildContext context, WidgetRef ref) async {
    try {
      // UserNotifier からユーザー情報を取得
      final user = ref.watch(userProvider);
      if (user == null) {
        throw Exception('ユーザー情報がロードされていません');
      }

      final internalUserId = user.id;
      if (internalUserId == null) {
        throw Exception('内部ユーザーIDが存在しません');
      }
      debugPrint('取得した内部ユーザーID: $internalUserId');

      // Google Sign-Inの初期化
      final GoogleSignIn googleSignIn = GoogleSignIn(
        scopes: [
          'https://www.googleapis.com/auth/calendar.readonly',
          'https://www.googleapis.com/auth/calendar',
        ],
        clientId: DefaultFirebaseOptions.ios.iosClientId,
        serverClientId:
            'XXXXXX.apps.googleusercontent.com', // WebクライアントID
      );

      // Googleサインインを実行
      final account = await googleSignIn.signIn();
      if (account == null) {
        debugPrint('Googleサインインがキャンセルされました');
        return;
      }

      // 認証コードの取得
      final authCode = account.serverAuthCode;
      if (authCode == null) {
        throw Exception('認証コードが取得できませんでした');
      }
      debugPrint('取得した認証コード: $authCode');

      // サーバーに認証コードを送信してリフレッシュトークンを取得・保存
      final googleCalendarService =
          GoogleCalendarService(apiClient: ref.read(apiClientProvider));
      await googleCalendarService.sendAuthCode(
          internalUserId.toString(), authCode);
      debugPrint('リフレッシュトークンの取得に成功しました');
    } catch (e) {
      debugPrint('エラーが発生しました: $e');
    }
  }
}