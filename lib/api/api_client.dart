// courtesy of Kazutaka Yoshinaga
extension GoogleCalendarApi on ApiClient {
  Future sendAuthCode(String userId, String authCode) async {
    final url = Uri.parse('$apiUrl/app/google-calendar/tokens');
    final headers = await apiHeaders();

    final body = {
      'userId': int.parse(userId),
      'authCode': authCode,
    };

    final response = await http.post(
      url,
      headers: headers,
      body: jsonEncode(body),
    );

    if (response.statusCode != 200) {
      throw Exception('Failed to send auth code: ${response.body}');
    }
  }
}