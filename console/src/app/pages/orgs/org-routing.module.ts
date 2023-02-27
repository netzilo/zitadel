import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';

import { OrgDetailComponent } from './org-detail/org-detail.component';

const routes: Routes = [
  //   {
  //     path: 'idp',
  //     children: [
  //       {
  //         path: 'create',
  //         loadChildren: () => import('src/app/modules/providers/provider-oidc/provider-oidc.module'),
  //         canActivate: [RoleGuard],
  //         data: {
  //           roles: ['org.idp.write'],
  //           serviceType: PolicyComponentServiceType.MGMT,
  //         },
  //       },
  //       {
  //         path: ':id',
  //         loadChildren: () => import('src/app/modules/idp/idp.module'),
  //         canActivate: [RoleGuard],
  //         data: {
  //           roles: ['org.idp.read'],
  //           serviceType: PolicyComponentServiceType.MGMT,
  //         },
  //       },
  //     ],
  //   },
  {
    path: 'members',
    loadChildren: () => import('./org-members/org-members.module'),
  },
  {
    path: '',
    component: OrgDetailComponent,
  },
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule],
})
export class OrgRoutingModule {}
