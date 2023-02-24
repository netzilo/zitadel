import { OverlayModule } from '@angular/cdk/overlay';
import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MatLegacyButtonModule as MatButtonModule } from '@angular/material/legacy-button';
import { MatIconModule } from '@angular/material/icon';
import { MatLegacyMenuModule as MatMenuModule } from '@angular/material/legacy-menu';
import { MatLegacyProgressSpinnerModule as MatProgressSpinnerModule } from '@angular/material/legacy-progress-spinner';
import { MatLegacyTooltipModule as MatTooltipModule } from '@angular/material/legacy-tooltip';
import { RouterModule } from '@angular/router';
import { TranslateModule } from '@ngx-translate/core';
import { HasRoleModule } from 'src/app/directives/has-role/has-role.module';
import { HasRolePipeModule } from 'src/app/pipes/has-role-pipe/has-role-pipe.module';

import { NavComponent } from './nav.component';
import OnboardingCardModule from '../onboarding-card/onboarding-card.module';
import { MatLegacyProgressBarModule } from '@angular/material/legacy-progress-bar';

@NgModule({
  declarations: [NavComponent],
  imports: [
    CommonModule,
    OnboardingCardModule,
    FormsModule,
    ReactiveFormsModule,
    TranslateModule,
    MatIconModule,
    RouterModule,
    MatTooltipModule,
    HasRolePipeModule,
    MatLegacyProgressBarModule,
    HasRoleModule,
    MatMenuModule,
    MatButtonModule,
    MatProgressSpinnerModule,
    OverlayModule,
  ],
  exports: [NavComponent],
})
export class NavModule {}
