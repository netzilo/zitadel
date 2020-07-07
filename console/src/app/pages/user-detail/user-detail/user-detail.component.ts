import { Location } from '@angular/common';
import { Component, OnDestroy, OnInit } from '@angular/core';
import { ActivatedRoute, Params } from '@angular/router';
import { TranslateService } from '@ngx-translate/core';
import { Subscription } from 'rxjs';
import { ChangeType } from 'src/app/modules/changes/changes.component';
import {
    Gender,
    NotificationType,
    UserEmail,
    UserPhone,
    UserProfile,
    UserState,
    UserView,
} from 'src/app/proto/generated/management_pb';
import { AuthUserService } from 'src/app/services/auth-user.service';
import { MgmtUserService } from 'src/app/services/mgmt-user.service';
import { ToastService } from 'src/app/services/toast.service';

@Component({
    selector: 'app-user-detail',
    templateUrl: './user-detail.component.html',
    styleUrls: ['./user-detail.component.scss'],
})
export class UserDetailComponent implements OnInit, OnDestroy {
    public user!: UserView.AsObject;
    public genders: Gender[] = [Gender.GENDER_MALE, Gender.GENDER_FEMALE, Gender.GENDER_DIVERSE];
    public languages: string[] = ['de', 'en'];

    public isMgmt: boolean = false;
    private subscription: Subscription = new Subscription();
    public emailEditState: boolean = false;
    public phoneEditState: boolean = false;

    public ChangeType: any = ChangeType;
    public loading: boolean = false;

    public UserState: any = UserState;
    public copied: string = '';

    constructor(
        public translate: TranslateService,
        private route: ActivatedRoute,
        private toast: ToastService,
        private mgmtUserService: MgmtUserService,
        private _location: Location,
        public authUserService: AuthUserService,
    ) { }

    public ngOnInit(): void {
        this.subscription = this.route.params.subscribe(params => {
            this.loading = true;
            this.getData(params).then(() => {
                this.loading = false;
            }).catch(error => {
                this.loading = false;
            });
        });
    }

    public ngOnDestroy(): void {
        this.subscription.unsubscribe();
    }

    public changeState(newState: UserState): void {
        if (newState === UserState.USERSTATE_ACTIVE) {
            this.mgmtUserService.ReactivateUser(this.user.id).then(() => {
                this.toast.showInfo('reactivated User');
            }).catch(error => {
                this.toast.showError(error.message);
            });
        } else if (newState === UserState.USERSTATE_INACTIVE) {
            this.mgmtUserService.DeactivateUser(this.user.id).then(() => {
                this.toast.showInfo('deactivated User');
            }).catch(error => {
                this.toast.showError(error.message);
            });
        }
    }

    public saveProfile(profileData: UserProfile.AsObject): void {
        this.user.firstName = profileData.firstName;
        this.user.lastName = profileData.lastName;
        this.user.nickName = profileData.nickName;
        this.user.displayName = profileData.displayName;
        this.user.gender = profileData.gender;
        this.user.preferredLanguage = profileData.preferredLanguage;
        this.mgmtUserService
            .SaveUserProfile(
                this.user.id,
                this.user.firstName,
                this.user.lastName,
                this.user.nickName,
                this.user.preferredLanguage,
                this.user.gender)
            .then((data: UserProfile) => {
                this.toast.showInfo('Saved Profile');
                this.user = Object.assign(this.user, data.toObject());
            })
            .catch(data => {
                this.toast.showError(data.message);
            });
    }

    public resendVerification(): void {
        this.mgmtUserService.ResendEmailVerification(this.user.id).then(() => {
            this.toast.showInfo('Email was successfully sent!');
        }).catch(data => {
            this.toast.showError(data.message);
        });
    }

    public resendPhoneVerification(): void {
        this.mgmtUserService.ResendPhoneVerification(this.user.id).then(() => {
            this.toast.showInfo('Phoneverification was successfully sent!');
        }).catch(data => {
            this.toast.showError(data.message);
        });
    }

    public deletePhone(): void {
        this.mgmtUserService.RemoveUserPhone(this.user.id).then(() => {
            this.toast.showInfo('Phone removed with success!');
            this.user.phone = '';
            this.phoneEditState = false;
        }).catch(data => {
            this.toast.showError(data.message);
        });
    }

    public saveEmail(): void {
        this.emailEditState = false;
        this.mgmtUserService
            .SaveUserEmail(this.user.id, this.user.email).then((data: UserEmail) => {
                this.toast.showInfo('Saved Email');
                this.user.email = data.toObject().email;
            }).catch(data => {
                this.toast.showError(data.message);
            });
    }

    public savePhone(): void {
        this.phoneEditState = false;
        this.mgmtUserService
            .SaveUserPhone(this.user.id, this.user.phone).then((data: UserPhone) => {
                this.toast.showInfo('Saved Phone');
                this.user.phone = data.toObject().phone;
                this.phoneEditState = false;
            }).catch(data => {
                this.toast.showError(data.message);
            });
    }

    public navigateBack(): void {
        this._location.back();
    }

    private async getData({ id }: Params): Promise<void> {
        this.isMgmt = true;
        this.mgmtUserService.GetUserByID(id).then(user => {
            this.user = user.toObject();
        }).catch(err => {
            console.error(err);
        });
    }

    public sendSetPasswordNotification(): void {
        this.mgmtUserService.SendSetPasswordNotification(this.user.id, NotificationType.NOTIFICATIONTYPE_EMAIL)
            .then(() => {
                this.toast.showInfo('Set initial Password');
            }).catch(data => {
                this.toast.showError(data.message);
            });
    }

    public copytoclipboard(value: string): void {
        const selBox = document.createElement('textarea');
        selBox.style.position = 'fixed';
        selBox.style.left = '0';
        selBox.style.top = '0';
        selBox.style.opacity = '0';
        selBox.value = value;
        document.body.appendChild(selBox);
        selBox.focus();
        selBox.select();
        document.execCommand('copy');
        document.body.removeChild(selBox);
        this.copied = value;
        setTimeout(() => {
            this.copied = '';
        }, 3000);
    }
}
