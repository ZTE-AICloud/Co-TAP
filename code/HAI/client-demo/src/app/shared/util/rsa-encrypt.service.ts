declare let encryptRSA;

export class RSAEncryptService {
  /**
   * 使用RSA加密，支持超过245byte长度，使用分段加密。
   * encryptRSA是encrypt.min.js提供的分段加密方法，在index.html中引入了encrypt.min.js
   * @param content 加密内容
   */
  public encrypt(content: string): string {
    if (content && encryptRSA) {
      return encryptRSA(content);
    }
    return '';
  }
}