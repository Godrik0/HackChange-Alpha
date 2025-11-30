import {Component, Input} from '@angular/core';
import {NzCardComponent} from "ng-zorro-antd/card";

@Component({
  selector: 'app-action-card',
  imports: [
    NzCardComponent
  ],
  templateUrl: './action-card.component.html',
  styleUrl: './action-card.component.less',
})
export class ActionCardComponent {
  @Input() recommendation!: string;
}
