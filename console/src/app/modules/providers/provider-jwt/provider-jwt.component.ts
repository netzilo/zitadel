import { COMMA, ENTER, SPACE } from '@angular/cdk/keycodes';
import { Location } from '@angular/common';
import { Component, Injector, OnDestroy, OnInit, Type } from '@angular/core';
import { AbstractControl, UntypedFormControl, UntypedFormGroup, Validators } from '@angular/forms';
import { ActivatedRoute, Params, Router } from '@angular/router';
import { Subscription } from 'rxjs';
import { take } from 'rxjs/operators';
import {
  AddJWTProviderRequest as AdminAddJWTProviderRequest,
  UpdateJWTProviderRequest as AdminUpdateJWTProviderRequest,
} from 'src/app/proto/generated/zitadel/admin_pb';
import { OIDCMappingField, Provider } from 'src/app/proto/generated/zitadel/idp_pb';
import {
  AddJWTProviderRequest as MgmtAddJWTProviderRequest,
  UpdateJWTProviderRequest as MgmtUpdateJWTProviderRequest,
} from 'src/app/proto/generated/zitadel/management_pb';
import { AdminService } from 'src/app/services/admin.service';
import { Breadcrumb, BreadcrumbService, BreadcrumbType } from 'src/app/services/breadcrumb.service';
import { ManagementService } from 'src/app/services/mgmt.service';
import { ToastService } from 'src/app/services/toast.service';

import { PolicyComponentServiceType } from '../../policies/policy-component-types.enum';

@Component({
  selector: 'cnsl-provider-jwt',
  templateUrl: './provider-jwt.component.html',
  styleUrls: ['./provider-jwt.component.scss'],
})
export class ProviderJWTComponent implements OnInit, OnDestroy {
  public serviceType: PolicyComponentServiceType = PolicyComponentServiceType.MGMT;
  private service!: ManagementService | AdminService;
  public readonly separatorKeysCodes: number[] = [ENTER, COMMA, SPACE];
  public mappingFields: OIDCMappingField[] = [];

  private subscription?: Subscription;
  public projectId: string = '';

  public jwtFormGroup!: UntypedFormGroup;
  public loading: boolean = false;

  public provider?: Provider.AsObject;

  constructor(
    private router: Router,
    private route: ActivatedRoute,
    private toast: ToastService,
    private injector: Injector,
    private _location: Location,
    breadcrumbService: BreadcrumbService,
  ) {
    this.jwtFormGroup = new UntypedFormGroup({
      jwtName: new UntypedFormControl('', [Validators.required]),
      jwtHeaderName: new UntypedFormControl('', [Validators.required]),
      jwtIssuer: new UntypedFormControl('', [Validators.required]),
      jwtEndpoint: new UntypedFormControl('', [Validators.required]),
      jwtKeysEndpoint: new UntypedFormControl('', [Validators.required]),
      jwtStylingType: new UntypedFormControl(0),
      jwtAutoRegister: new UntypedFormControl(false),
    });

    this.route.data.pipe(take(1)).subscribe((data) => {
      this.serviceType = data.serviceType;
      switch (this.serviceType) {
        case PolicyComponentServiceType.MGMT:
          this.service = this.injector.get(ManagementService as Type<ManagementService>);
          this.mappingFields = [
            OIDCMappingField.OIDC_MAPPING_FIELD_PREFERRED_USERNAME,
            OIDCMappingField.OIDC_MAPPING_FIELD_EMAIL,
          ];
          const bread: Breadcrumb = {
            type: BreadcrumbType.ORG,
            routerLink: ['/org'],
          };

          breadcrumbService.setBreadcrumb([bread]);
          break;
        case PolicyComponentServiceType.ADMIN:
          this.service = this.injector.get(AdminService as Type<AdminService>);
          this.mappingFields = [
            OIDCMappingField.OIDC_MAPPING_FIELD_PREFERRED_USERNAME,
            OIDCMappingField.OIDC_MAPPING_FIELD_EMAIL,
          ];

          const iamBread = new Breadcrumb({
            type: BreadcrumbType.ORG,
            name: 'Instance',
            routerLink: ['/instance'],
          });
          breadcrumbService.setBreadcrumb([iamBread]);
          break;
      }
    });
  }

  public ngOnInit(): void {
    this.subscription = this.route.params.subscribe((params) => this.getData(params));
  }

  public ngOnDestroy(): void {
    this.subscription?.unsubscribe();
  }

  private getData({ projectid }: Params): void {
    this.projectId = projectid;
  }

  public submitForm(): void {
    this.provider ? this.updateJWTProvider() : this.addJWTProvider();
  }

  public addJWTProvider(): void {
    if (this.serviceType === PolicyComponentServiceType.MGMT) {
      const req = new MgmtAddJWTProviderRequest();

      req.setName(this.jwtName?.value);
      req.setHeaderName(this.jwtHeaderName?.value);
      req.setIssuer(this.jwtIssuer?.value);
      req.setJwtEndpoint(this.jwtEndpoint?.value);
      req.setKeysEndpoint(this.jwtKeysEndpoint?.value);

      this.loading = true;
      //   (this.service as ManagementService)
      //     .addOrgJWTIDP(req)
      //     .then((idp) => {
      //       setTimeout(() => {
      //         this.loading = false;
      //         this.router.navigate([
      //           this.serviceType === PolicyComponentServiceType.MGMT
      //             ? 'org'
      //             : this.serviceType === PolicyComponentServiceType.ADMIN
      //             ? 'iam'
      //             : '',
      //           'policy',
      //           'login',
      //         ]);
      //       }, 2000);
      //     })
      //     .catch((error) => {
      //       this.toast.showError(error);
      //   this.loading=false;
      //     });
    } else if (PolicyComponentServiceType.ADMIN) {
      const req = new AdminAddJWTProviderRequest();

      req.setName(this.jwtName?.value);
      req.setHeaderName(this.jwtHeaderName?.value);
      req.setIssuer(this.jwtIssuer?.value);
      req.setJwtEndpoint(this.jwtEndpoint?.value);
      req.setKeysEndpoint(this.jwtKeysEndpoint?.value);

      this.loading = true;
      //   (this.service as AdminService)
      //     .addJWTIDP(req)
      //     .then((idp) => {
      //       setTimeout(() => {
      //         this.loading = false;
      //         this.router.navigate([
      //           this.serviceType === PolicyComponentServiceType.MGMT
      //             ? 'org'
      //             : this.serviceType === PolicyComponentServiceType.ADMIN
      //             ? 'iam'
      //             : '',
      //           'policy',
      //           'login',
      //         ]);
      //       }, 2000);
      //     })
      //     .catch((error) => {
      //       this.toast.showError(error);
      //   this.loading=false;
      //     });
    }
  }

  public updateJWTProvider(): void {
    if (this.provider) {
      if (this.serviceType === PolicyComponentServiceType.MGMT) {
        const req = new MgmtUpdateJWTProviderRequest();
        req.setId(this.provider.id);
        req.setName(this.jwtName?.value);
        req.setHeaderName(this.jwtHeaderName?.value);
        req.setIssuer(this.jwtIssuer?.value);
        req.setJwtEndpoint(this.jwtEndpoint?.value);
        req.setKeysEndpoint(this.jwtKeysEndpoint?.value);

        this.loading = true;
        //   (this.service as ManagementService)
        //     .addOrgJWTIDP(req)
        //     .then((idp) => {
        //       setTimeout(() => {
        //         this.loading = false;
        //         this.router.navigate([
        //           this.serviceType === PolicyComponentServiceType.MGMT
        //             ? 'org'
        //             : this.serviceType === PolicyComponentServiceType.ADMIN
        //             ? 'iam'
        //             : '',
        //           'policy',
        //           'login',
        //         ]);
        //       }, 2000);
        //     })
        //     .catch((error) => {
        //       this.toast.showError(error);
        // this.loading=false;
        //     });
      } else if (PolicyComponentServiceType.ADMIN) {
        const req = new AdminUpdateJWTProviderRequest();
        req.setId(this.provider.id);
        req.setName(this.jwtName?.value);
        req.setHeaderName(this.jwtHeaderName?.value);
        req.setIssuer(this.jwtIssuer?.value);
        req.setJwtEndpoint(this.jwtEndpoint?.value);
        req.setKeysEndpoint(this.jwtKeysEndpoint?.value);

        this.loading = true;
        //   (this.service as AdminService)
        //     .addJWTIDP(req)
        //     .then((idp) => {
        //       setTimeout(() => {
        //         this.loading = false;
        //         this.router.navigate([
        //           this.serviceType === PolicyComponentServiceType.MGMT
        //             ? 'org'
        //             : this.serviceType === PolicyComponentServiceType.ADMIN
        //             ? 'iam'
        //             : '',
        //           'policy',
        //           'login',
        //         ]);
        //       }, 2000);
        //     })
        //     .catch((error) => {
        //       this.toast.showError(error);
        // this.loading=false;
        //     });
      }
    }
  }

  public close(): void {
    this._location.back();
  }

  public get jwtName(): AbstractControl | null {
    return this.jwtFormGroup.get('jwtName');
  }

  public get jwtHeaderName(): AbstractControl | null {
    return this.jwtFormGroup.get('jwtHeaderName');
  }

  public get jwtIssuer(): AbstractControl | null {
    return this.jwtFormGroup.get('jwtIssuer');
  }

  public get jwtEndpoint(): AbstractControl | null {
    return this.jwtFormGroup.get('jwtEndpoint');
  }

  public get jwtKeysEndpoint(): AbstractControl | null {
    return this.jwtFormGroup.get('jwtKeysEndpoint');
  }
}
