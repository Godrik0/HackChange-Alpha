import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import {RouterModule} from "@angular/router";
import {ClientMetricsComponent} from "@features/client-metrics/pages/client-metrics/client-metrics.component";
import {IncomeChartComponent} from "@features/client-metrics/components/income-chart/income-chart.component";



@NgModule({
  declarations: [],
  imports: [
    CommonModule,
    RouterModule.forChild([
      {path: 'client/:id', component: ClientMetricsComponent},
    ]),
    ClientMetricsComponent,
    IncomeChartComponent
  ]
})
export class ClientMetricsModule { }
