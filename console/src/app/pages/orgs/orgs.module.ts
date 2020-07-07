import { CommonModule } from '@angular/common';
import { NgModule, NO_ERRORS_SCHEMA } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatButtonToggleModule } from '@angular/material/button-toggle';
import { MatDialogModule } from '@angular/material/dialog';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatIconModule } from '@angular/material/icon';
import { MatInputModule } from '@angular/material/input';
import { MatMenuModule } from '@angular/material/menu';
import { MatTabsModule } from '@angular/material/tabs';
import { MatTooltipModule } from '@angular/material/tooltip';
import { TranslateModule } from '@ngx-translate/core';
import { HasRoleModule } from 'src/app/directives/has-role/has-role.module';
import { CardModule } from 'src/app/modules/card/card.module';
import { MetaLayoutModule } from 'src/app/modules/meta-layout/meta-layout.module';
import { WarnDialogModule } from 'src/app/modules/warn-dialog/warn-dialog.module';

import { ChangesModule } from '../../modules/changes/changes.module';
import { OrgContributorsModule } from './org-contributors/org-contributors.module';
import { OrgDetailComponent } from './org-detail/org-detail.component';
import { OrgGridComponent } from './org-grid/org-grid.component';
import { OrgsRoutingModule } from './orgs-routing.module';
import { PolicyGridComponent } from './policy-grid/policy-grid.component';

@NgModule({
    declarations: [OrgDetailComponent, OrgGridComponent, PolicyGridComponent],
    imports: [
        CommonModule,
        OrgsRoutingModule,
        OrgContributorsModule,
        WarnDialogModule,
        FormsModule,
        HasRoleModule,
        MatFormFieldModule,
        MatInputModule,
        MatButtonModule,
        MatDialogModule,
        CardModule,
        MatIconModule,
        ReactiveFormsModule,
        MatButtonToggleModule,
        MetaLayoutModule,
        MatTabsModule,
        MatTooltipModule,
        MatMenuModule,
        ChangesModule,
        TranslateModule,
    ],
    exports: [],
    schemas: [NO_ERRORS_SCHEMA],
})
export class OrgsModule { }
