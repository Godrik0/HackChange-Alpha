import {Component, OnInit} from '@angular/core';
import {NzContentComponent, NzHeaderComponent, NzLayoutComponent} from "ng-zorro-antd/layout";
import {NzCardComponent} from "ng-zorro-antd/card";
import {NzColDirective, NzRowDirective} from "ng-zorro-antd/grid";
import {AsyncPipe, CommonModule, NgForOf, NgOptimizedImage} from "@angular/common";
import {IncomeChartComponent} from "@features/client-metrics/components/income-chart/income-chart.component";
import {CreditSliderComponent} from "@features/client-metrics/components/credit-slider/credit-slider.component";
import {BulletPointsComponent} from "@features/client-metrics/components/bullet-points/bullet-points.component";
import {ActionCardComponent} from "@features/client-metrics/components/action-card/action-card.component";
import {FormsModule} from "@angular/forms";
import {ActivatedRoute, RouterLink} from "@angular/router";
import {getFullName} from "@core/utils/utils";
import {Client} from "@core/models/client";
import {DataService} from "@core/services/data.service";
import {NzEmptyComponent} from "ng-zorro-antd/empty";
import {Observable} from "rxjs";
import {Scoring} from "@core/models/scoring";

@Component({
  selector: 'app-client-metrics',
  imports: [
    CommonModule,
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
    RouterLink,
    NzEmptyComponent,
    AsyncPipe,
  ],
  templateUrl: './client-metrics.component.html',
  styleUrl: './client-metrics.component.less',
})
export class ClientMetricsComponent implements OnInit {
  clientID: number | null = null;
  client$: Observable<Client> | undefined;
  scoring: Scoring | undefined;

  constructor(private route: ActivatedRoute, private dataService: DataService) {
  }

  ngOnInit() {
    this.clientID = Number(this.route.snapshot.paramMap.get('id')) ?? null;
    this.client$ = this.dataService.getClient(this.clientID);
    this.dataService.getScoringClient(this.clientID).subscribe(
      (scoring) => {
        this.scoring = scoring;
      }
    )
  }

  get incomeData() {
    return [
      { name: 'Текущий доход', value: this.scoring?.income },
      { name: 'Прогнозируемый доход', value: this.scoring?.predict_income },
    ]
  }

  get recommendations() {
    return this.scoring?.recommendations;
  }

  get selectedLimit() {
    return this.scoring?.credit_limit ?? 0;
  }

  get bulletPoints() {
    return {
      "negative_factors": this.scoring?.negative_factors!,
      "positive_factors": this.scoring?.positive_factors!
    }
  }

  protected readonly getFullName = getFullName;
}
