import {Component, Input} from '@angular/core';
import {NzSplitterModule} from "ng-zorro-antd/splitter";
import {NzFlexDirective} from "ng-zorro-antd/flex";
import {NgForOf} from "@angular/common";

@Component({
  selector: 'app-bullet-points',
  imports: [
    NzSplitterModule,
    NzFlexDirective,
    NgForOf
  ],
  templateUrl: './bullet-points.component.html',
  styleUrl: './bullet-points.component.less',
})
export class BulletPointsComponent {
  @Input() points: string[] = [];
}
