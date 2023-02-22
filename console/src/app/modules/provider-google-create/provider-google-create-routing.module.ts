import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';

import { ProviderGoogleCreateComponent } from './provider-google-create.component';

const routes: Routes = [
  {
    path: '',
    component: ProviderGoogleCreateComponent,
    data: { animation: 'DetailPage' },
  },
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule],
})
export class ProviderGoogleCreateRoutingModule {}
