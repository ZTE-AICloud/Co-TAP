import { RouterModule, Routes } from '@angular/router';
import { TaskComponent } from './task/task.component';

const routes: Routes = [
  {
    path: '',
    redirectTo: 'task',
    pathMatch: 'full'
  },
  {
    path: 'task',
    component: TaskComponent,
  },
];

export const PageRoutingModule = RouterModule.forChild(routes);