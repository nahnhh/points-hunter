// courtesy of Kazutaka Yoshinaga
import 'package:clocky_app/api/api_client.dart';

class GoogleCalendarService {
  final ApiClient apiClient;
  GoogleCalendarService({required this.apiClient});

  // 認可コードをバックエンドに送信
  Future sendAuthCode(String userId, String authCode) async {
    await apiClient.sendAuthCode(userId, authCode);
  }
}