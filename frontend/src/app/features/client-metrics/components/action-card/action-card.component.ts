import {Component, Input} from '@angular/core';
import {NzButtonComponent} from "ng-zorro-antd/button";
import {NzCardComponent} from "ng-zorro-antd/card";

@Component({
  selector: 'app-action-card',
  imports: [
    NzButtonComponent,
    NzCardComponent
  ],
  templateUrl: './action-card.component.html',
  styleUrl: './action-card.component.less',
})
export class ActionCardComponent {
  @Input() data: any;
}
