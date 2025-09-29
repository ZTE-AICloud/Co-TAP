export class RestAPI {
  private static readonly baseUrl = '/api/v1/zui/ai_dome';
  public static readonly chatUrl = `${this.baseUrl}/chat`;
  public static readonly taskUrl = `${this.baseUrl}/task`;
}