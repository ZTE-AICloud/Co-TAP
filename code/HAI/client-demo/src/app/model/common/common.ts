export interface HttpErrorMessage {
  zh: string;
  en: string;
}

export interface HttpResponse {
  items: any;
  message: HttpErrorMessage;
  errCode: string;
}