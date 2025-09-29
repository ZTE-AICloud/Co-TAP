import { Component, NgZone, ViewContainerRef } from '@angular/core';
import { Router } from '@angular/router';
import { TranslateService } from '@ngx-translate/core';
import { PlxI18nService } from 'paletx/plx-i18n';
import { PlxMessage } from 'paletx/plx-message';
import { CommonUtil } from './shared/util/common-util';
import { Constant } from './shared/util/constant';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html'
})
export class AppComponent {
  constructor(
    private _ngZone: NgZone,
    private routing: Router,
    private vRef: ViewContainerRef,
    private plxI18nService: PlxI18nService,
    private translate: TranslateService,
    private plxMessage: PlxMessage,
  ) {
    // 设置菜单跳转需要引用的对象，用于菜单跳转传递参数
    window['angularComponent'] = { component: this, zone: _ngZone };
    // 设置plxMessage组件全局提示
    this.plxMessage.setRootViewContainerRef(this.vRef);
    this.translate.addLangs(Constant.SUPPORT_LANGUAGES);
    this.translate.setDefaultLang('en-US');

    const lang = CommonUtil.getValidLanguage();
    this.translate.use(lang);
    this.plxI18nService.setLocale(lang);
  }
}