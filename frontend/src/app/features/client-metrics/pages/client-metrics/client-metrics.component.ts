import { Component } from '@angular/core';
import {NzContentComponent, NzHeaderComponent, NzLayoutComponent} from "ng-zorro-antd/layout";
import {NzCardComponent} from "ng-zorro-antd/card";
import {NzColDirective, NzRowDirective} from "ng-zorro-antd/grid";
import {NgForOf, NgOptimizedImage} from "@angular/common";
import {IncomeChartComponent} from "@features/client-metrics/components/income-chart/income-chart.component";
import {CreditSliderComponent} from "@features/client-metrics/components/credit-slider/credit-slider.component";
import {BulletPointsComponent} from "@features/client-metrics/components/bullet-points/bullet-points.component";
import {ActionCardComponent} from "@features/client-metrics/components/action-card/action-card.component";
import {FormsModule} from "@angular/forms";
import {RouterLink} from "@angular/router";

@Component({
  selector: 'app-client-metrics',
  imports: [
    NzLayoutComponent,
    NzHeaderComponent,
    NzContentComponent,
    NzCardComponent,
    NzRowDirective,
    NgForOf,
    NzColDirective,
    IncomeChartComponent,
    CreditSliderComponent,
    BulletPointsComponent,
    ActionCardComponent,
    FormsModule,
    NgOptimizedImage,
    RouterLink
  ],
  templateUrl: './client-metrics.component.html',
  styleUrl: './client-metrics.component.less',
})
export class ClientMetricsComponent {
  selectedLimit = 40000;
  incomeData = [
    { x: '32 000', y: 1 },
    { x: '55 000', y: 2 }
  ];
  bulletPoints = ['высокий рост дохода', 'нет просрочек по кредитам'];
  actionCards = [{ title: 'Отправить 1' }, { title: 'Отправить 2' }, { title: 'Отправить 3' }, { title: 'Отправить 4' }]
}
