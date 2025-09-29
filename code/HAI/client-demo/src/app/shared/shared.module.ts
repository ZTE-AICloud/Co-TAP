import { CommonModule } from '@angular/common';
import { ModuleWithProviders, NgModule } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { PaletXModule } from './paletx.module';
import { RSAEncryptService } from './util/rsa-encrypt.service';

@NgModule({
  declarations: [],
  imports: [
    CommonModule,
    FormsModule,
    ReactiveFormsModule,
    PaletXModule,
  ],
  exports: [
    CommonModule,
    FormsModule,
    ReactiveFormsModule,
    PaletXModule,
  ],
})
export class SharedModule {
  static forRoot(): ModuleWithProviders<SharedModule> {
    return {
      ngModule: SharedModule,
      providers: [
        RSAEncryptService,
      ]
    };
  }
}