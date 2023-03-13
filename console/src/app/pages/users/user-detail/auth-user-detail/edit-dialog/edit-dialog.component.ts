import { KeyValue, KeyValuePipe } from '@angular/common';
import { Component, Inject, OnInit } from '@angular/core';
import { AbstractControl, FormControl, FormGroup, UntypedFormBuilder, UntypedFormControl, UntypedFormGroup, Validators } from '@angular/forms';
import {
  MatLegacyDialogRef as MatDialogRef,
  MAT_LEGACY_DIALOG_DATA as MAT_DIALOG_DATA,
} from '@angular/material/legacy-dialog';
import { TranslateService } from '@ngx-translate/core';
import { requiredValidator } from 'src/app/modules/form-field/validators/validators';
import { CountryCallingCodesService, CountryPhoneCode } from 'src/app/services/country-calling-codes.service';
import { formatPhone } from 'src/app/utils/formatPhone';

export enum EditDialogType {
  PHONE = 1,
  EMAIL = 2,
}

@Component({
  selector: 'cnsl-edit-dialog',
  templateUrl: './edit-dialog.component.html',
  styleUrls: ['./edit-dialog.component.scss'],
})
export class EditDialogComponent implements OnInit {
  public controlKey = "editingField"
  public isPhone: boolean = false;
  public isVerified: boolean = false;
  public phoneCountry: string = 'CH';
  public dialogForm!: UntypedFormGroup;
  public EditDialogType: any = EditDialogType;
  public selected: CountryPhoneCode | undefined;
  public countryPhoneCodes: CountryPhoneCode[] = [];
  constructor(
    public dialogRef: MatDialogRef<EditDialogComponent>,
    @Inject(MAT_DIALOG_DATA) public data: any,
    private countryCallingCodesService: CountryCallingCodesService,
    private translateSvc: TranslateService,
    private kvPipe: KeyValuePipe
  ) {
    if (data.type === EditDialogType.PHONE) {
      this.isPhone = true;
    }
    this.dialogForm = new FormGroup({
      [this.controlKey]: new UntypedFormControl(data.value, data.validator || requiredValidator)
    })
  }

  public setCountryCallingCode(): void {
    console.log(this)
    let value = (this.dialogForm.controls[this.controlKey]?.value as string) || '';
    this.countryPhoneCodes.forEach(code => value = value.replace(`+${code.countryCallingCode}`, ""))
    value = value.trim()
    this.dialogForm.controls[this.controlKey]?.setValue('+' + this.selected?.countryCallingCode + ' ' + value);
    console.log(this)
  }

  ngOnInit(): void {
    if (this.isPhone) {
      // Get country phone codes and set selected flag to guessed country or default country
      this.countryPhoneCodes = this.countryCallingCodesService.getCountryCallingCodes();
      const phoneNumber = formatPhone(this.dialogForm.controls[this.controlKey]?.value);
      this.selected = this.countryPhoneCodes.find((code) => code.countryCode === phoneNumber.country);
      this.dialogForm.controls[this.controlKey].setValue(phoneNumber.phone);
    }
  }

  closeDialog(): void {
    this.dialogRef.close();
  }

  closeDialogWithValue(): void {
    this.dialogRef.close({ value: this.dialogForm.controls[this.controlKey].value, isVerified: this.isVerified });
  }

  public get ctrl(): AbstractControl {
    return this.dialogForm.get(this.controlKey) as AbstractControl;
  }
}

