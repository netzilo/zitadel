import { Component, EventEmitter, Input, Output } from '@angular/core';
import { ManagementService } from '../../../services/mgmt.service';
import { AdminService } from '../../../services/admin.service';
import { PolicyComponentServiceType } from '../../policies/policy-component-types.enum';

export interface CopyUrl {
  label: string;
  url: string;
  downloadable?: boolean;
}

@Component({
  selector: 'cnsl-provider-next',
  templateUrl: './provider-next.component.html',
  styleUrls: ['./provider-next.component.scss'],
})
export class ProviderNextComponent {
  @Input() copyUrls?: CopyUrl[] | null;
  @Input() autofillLink?: string | null;
  @Input() activateLink?: string | null;
  @Input() configureProvider?: boolean | null;
  @Input() configureTitle?: string;
  @Input() configureDescription?: string;
  @Input() configureLink?: string;
  @Input() expanded?: boolean;
  @Output() activate = new EventEmitter<void>();

  constructor() {}
}
