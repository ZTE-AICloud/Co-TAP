import { ValidatorFn, Validators } from '@angular/forms';
import { Constant } from './constant';

export class CommonUtil {
  private static getStandardLanguage(): string {
    // 获取uportal框架的语言信息，tools.js提供getLanguage()方法
    let uportalLang = '';
    if (window['getLanguage']) {
      uportalLang = window['getLanguage']();
    }
    return uportalLang || window.navigator.language;
  }

  public static getValidLanguage(): string {
    const uportalLang = this.getStandardLanguage();
    const validLang = Constant.SUPPORT_LANGUAGES.includes(uportalLang) ? uportalLang : 'en-US';
    return validLang;
  }

  /* Started by AICoder, pid:sea699ff13794331482808243052f4087fd8b15e */
  /**
   * 检查当前语言是否为中文
   * @param currentLang 当前语言
   * @returns 当前语言是否为中文
   */
  public static isZhLanguage(currentLang: string): boolean {
    return currentLang === 'zh-CN';
  }
  /* Ended by AICoder, pid:sea699ff13794331482808243052f4087fd8b15e */

  /* Started by AICoder, pid:t578f6c3b0f83ea1428409ba5077570fe0c32b18 */
  public static getRequiredValidator(value: string | null | undefined): ValidatorFn[] {
    return value ? [] : [Validators.required];
  }
  /* Ended by AICoder, pid:t578f6c3b0f83ea1428409ba5077570fe0c32b18 */

  /* Started by AICoder, pid:7419dv39d6pd303148ad0873f0a44c197fa9b035 */
  public static removeEmptyValues(obj: any): any {
    if (typeof obj !== 'object' || obj === null) {
      return obj;
    }
    if (Array.isArray(obj)) {
      return obj;
    }

    const result = {};
    for (const key in obj) {
      if (Object.prototype.hasOwnProperty.call(obj, key)) {
        const value = obj[key];
        if (value !== null && value !== undefined && value !== '') {
          result[key] = value;
        }
      }
    }
    return result;
  }
  /* Ended by AICoder, pid:7419dv39d6pd303148ad0873f0a44c197fa9b035 */
}