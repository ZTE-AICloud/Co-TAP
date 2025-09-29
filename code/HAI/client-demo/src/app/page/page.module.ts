import { NgModule } from '@angular/core';
import { SharedModule } from '../shared/shared.module';
import { PageRoutingModule } from './page-routing.module';
import { TaskComponent } from './task/task.component';

@NgModule({
  declarations: [
    TaskComponent,
  ],
  imports: [
    SharedModule,
    PageRoutingModule,
  ]
})
export class PageModule { }