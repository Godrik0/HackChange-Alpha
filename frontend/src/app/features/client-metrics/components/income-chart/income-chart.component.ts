// income-chart.component.ts
import {Component, Input} from '@angular/core';
import {BarChartModule, Color, ScaleType} from '@swimlane/ngx-charts';

@Component({
  selector: 'app-income-chart',
  templateUrl: './income-chart.component.html',
  imports: [
    BarChartModule
  ]
})
export class IncomeChartComponent {
  @Input() data = {};

  // data = [
  //   { name: 'Jan', value: 100 },
  //   { name: 'Feb', value: 150 },
  //   { name: 'Mar', value: 80 }
  // ];

  view: [number, number] = [700, 400];
  showXAxis = true;
  showYAxis = true;
  showLegend = false;
  showXAxisLabel = true;
  showYAxisLabel = true;
  xAxisLabel = 'Месяц';
  yAxisLabel = 'Значение';
  gradient = false;

  // Красная схема
  redColorScheme: Color = {
    name: 'red',
    selectable: true,
    group: ScaleType.Ordinal,
    domain: ['#ff4d4f'] // красный цвет столбцов
  };
}
