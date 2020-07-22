import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatTooltipModule } from '@angular/material/tooltip';
import { TranslateModule } from '@ngx-translate/core';

import { RefreshTableComponent } from './refresh-table.component';



@NgModule({
    declarations: [RefreshTableComponent],
    imports: [
        CommonModule,
        MatButtonModule,
        MatIconModule,
        TranslateModule,
        FormsModule,
        MatTooltipModule,
        MatProgressSpinnerModule,
    ],
    exports: [
        RefreshTableComponent,
    ],
})
export class RefreshTableModule { }
