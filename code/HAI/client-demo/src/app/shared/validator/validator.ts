import { AbstractControl } from '@angular/forms';

const NAME_REG = /^[0-9a-zA-Z\u4e00-\u9fa5]+([.0-9a-zA-Z\u4e00-\u9fa5_-]*[0-9a-zA-Z\u4e00-\u9fa5]$)?/;
export const IPV4_REG = /^((0|[1-9]\d?|1\d\d|2[0-4]\d|25[0-5])\.){3}(0|[1-9]\d?|1\d\d|2[0-4]\d|25[0-5])$/;
export const IPV6_REG = /^\s*((([0-9A-Fa-f]{1,4}:){7}([0-9A-Fa-f]{1,4}|:))|(([0-9A-Fa-f]{1,4}:){6}(:[0-9A-Fa-f]{1,4}|((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3})|:))|(([0-9A-Fa-f]{1,4}:){5}(((:[0-9A-Fa-f]{1,4}){1,2})|:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3})|:))|(([0-9A-Fa-f]{1,4}:){4}(((:[0-9A-Fa-f]{1,4}){1,3})|((:[0-9A-Fa-f]{1,4})?:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(([0-9A-Fa-f]{1,4}:){3}(((:[0-9A-Fa-f]{1,4}){1,4})|((:[0-9A-Fa-f]{1,4}){0,2}:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(([0-9A-Fa-f]{1,4}:){2}(((:[0-9A-Fa-f]{1,4}){1,5})|((:[0-9A-Fa-f]{1,4}){0,3}:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(([0-9A-Fa-f]{1,4}:){1}(((:[0-9A-Fa-f]{1,4}){1,6})|((:[0-9A-Fa-f]{1,4}){0,4}:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(:(((:[0-9A-Fa-f]{1,4}){1,7})|((:[0-9A-Fa-f]{1,4}){0,5}:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:)))(%.+)?\s*$/;

export const validateName = (c: AbstractControl): { [key: string]: any } => {
  if (!c.value) {
    return null;
  }
  const result = c.value.match(NAME_REG);
  if (!result || !result.length || result[0] !== c.value) {
    return { name: true };
  }
  return null;
};

export const validateIP = (c: AbstractControl): { [key: string]: any } => {
  if (!c.value) {
    return null;
  }

  const v = c.value;
  if (!v.match(IPV4_REG) && !v.match(IPV6_REG)) {
    return { ip: true };
  }
  return null;
};

export const validateDir = (c: AbstractControl): { [key: string]: any } => {
  if (!c.value) {
    return null;
  }

  const v = c.value;
  // 匹配只输入根目录/，或者是以/开头的目录
  if (!/(^\/$)|(^\/[\w]+(?:[\/][\w]+)*$)/.test(v)) {
    return { dir: true };
  } else if (v.length > 64) {
    return { dir: true };
  }
  return null;
};

export const validateDesc = (c: AbstractControl): { [key: string]: any } => {
  if (!c.value) {
    return null;
  }

  const v = c.value;
  if (/^[\+\-=@\s]+/.test(v)) {
    return { desc: true };
  }
  return null;
};