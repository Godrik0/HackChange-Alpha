import {Component, Input} from '@angular/core';
import {FormsModule} from "@angular/forms";
import {NzMarks, NzSliderComponent} from "ng-zorro-antd/slider";

@Component({
  selector: 'app-credit-slider',
  imports: [
    FormsModule,
    NzSliderComponent
  ],
  templateUrl: './credit-slider.component.html',
  styleUrl: './credit-slider.component.less',
})
export class CreditSliderComponent {
  min = 0;
  @Input() max = 100000;
  @Input() optimal = 70000;

  get marks(): NzMarks {
    return {
      0: '0 ₽',
      [this.max]: `${this.max.toLocaleString('ru-RU')} ₽`,
    };
  }
}
