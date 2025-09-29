import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { TranslateService } from '@ngx-translate/core';

@Injectable()
export class HttpErrorService {
  constructor(
    private translate: TranslateService,
    private http: HttpClient,
  ) { }

  /* Started by AICoder, pid:i27a024969s398e14600083f6095313e6ab3bc0b */
  public getErrorMessage(error: any): string {
    const SERVER_ERROR_KEY = 'common.message.serverError';
    const NO_RIGHT_KEY = 'common.message.noRight';

    // 处理 501 状态码情况
    if (error.status === 501) {
      this.reloadPage();
      return this.translate.instant(SERVER_ERROR_KEY);
    }

    // 统一处理非 501 情况
    if (!error.error) {
      return this.translate.instant(SERVER_ERROR_KEY);
    }

    // 解析错误信息
    let err: any;
    try {
      err = JSON.parse(error.error);
    } catch {
      err = error.error;
    }
    // code用于otaf场景，errorCode用于idos场景
    const { labels, message, code, errorCode } = err;

    // 处理特定错误类型
    if (this.isAuthorizationError(message, code, errorCode)) {
      return this.translate.instant(NO_RIGHT_KEY);
    }

    // 处理多语言标签
    const localizedLabel = this.getLocalizedLabel(labels);
    if (localizedLabel) {
      return localizedLabel;
    }

    return this.translate.instant(SERVER_ERROR_KEY);
  }
  /* Ended by AICoder, pid:i27a024969s398e14600083f6095313e6ab3bc0b */

  /* Started by AICoder, pid:5917bpc88br9a9914deb0af48030632d2ef5177b */
  private getLocalizedLabel(labels?: Record<string, string>): string | null {
    if (!labels) {
      return null;
    }

    let serverKey: string;
    let platformKey: string;

    switch (this.translate.currentLang) {
      case 'zh-CN':
        serverKey = 'zh';
        platformKey = 'zh_CN';
        break;
      case 'es-ES':
        serverKey = 'es';
        platformKey = 'es_ES';
        break;
      default:
        serverKey = 'en';
        platformKey = 'en_US';
        break;
    }

    return labels[serverKey] || labels[platformKey] || null;
  }
  /* Ended by AICoder, pid:5917bpc88br9a9914deb0af48030632d2ef5177b */

  /* Started by AICoder, pid:u27a0r4969e398e14600083f6095311e6ab8bc0b */
  private isAuthorizationError(message: string, code: number, errorCode: number): boolean {
    // 403的情况，备机环境和sm鉴权失败都会返回message
    // sm鉴权失败，对message信息做国际化处理
    const AUTH_MESSAGES = [
      'extent request',
      'no authorized access',
      'this url not authority'
    ];
    const CODE_NO_RIGHT = 3329001;
    const ERROR_CODE_NO_RIGHT = 65990002;

    return AUTH_MESSAGES.includes(message) || code === CODE_NO_RIGHT || errorCode === ERROR_CODE_NO_RIGHT;
  }
  /* Ended by AICoder, pid:u27a0r4969e398e14600083f6095311e6ab8bc0b */

  /**
   * @description 当服务端返回501的时候，需要判断下前端容器是否正常运行，去请求国际化文件，如果前端容器不正常了，则需要重新加载下资源
   * 因为在chrome下，当ssm服务被停了，点击左侧菜单重新加载ssm页面，界面显示501
   * 在firefox下，当ssm服务被停了，点击左侧菜单重新加载ssm页面，界面存在缓存，会加载出ssm界面，但是没有国际化，存在问题，故需要重新加载资源
   */
  private reloadPage(): void {
    this.http.get('./assets/i18n/zh-CN.properties').subscribe({
      next: () => {
        // 无需处理
      },
      error: error => {
        if (error.status === 501) {
          window.location.reload();
        }
      }
    });
  }
}