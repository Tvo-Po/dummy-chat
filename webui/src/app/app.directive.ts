import { Directive } from '@angular/core';
import {
  AbstractControl,
  NG_VALIDATORS,
  Validator,
  ValidationErrors,
  ValidatorFn
} from '@angular/forms';

export function forbiddenUsernameValidator(): ValidatorFn {
  return (control: AbstractControl): ValidationErrors | null => {
    if (control.value === "_system_") {
      return {forbiddenUsername: "Name unavailable"}
    }
    return null;
  };
}

@Directive({
  selector: '[forbiddenUsername]',
  providers: [
    {
      provide: NG_VALIDATORS,
      useExisting: ForbiddenValidatorDirective,
      multi: true,
    },
  ],
  standalone: true,
})
export class ForbiddenValidatorDirective implements Validator {
  validate(control: AbstractControl): ValidationErrors | null {
    return forbiddenUsernameValidator()(control);
  }
}
