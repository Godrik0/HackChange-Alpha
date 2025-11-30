import {Component, OnInit, OnDestroy, PLATFORM_ID, Inject, ChangeDetectorRef} from '@angular/core';
import {NzContentComponent, NzHeaderComponent, NzLayoutComponent} from "ng-zorro-antd/layout";
import {NzCardComponent} from "ng-zorro-antd/card";
import {NzColDirective, NzRowDirective} from "ng-zorro-antd/grid";
import {CommonModule, NgForOf, NgOptimizedImage, isPlatformBrowser} from "@angular/common";
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
import {Scoring} from "@core/models/scoring";
import {Subject, takeUntil, switchMap} from "rxjs";

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
  ],
  templateUrl: './client-metrics.component.html',
  styleUrl: './client-metrics.component.less',
})
export class ClientMetricsComponent implements OnInit, OnDestroy {
  clientID: number | null = null;
  scoring: Scoring | undefined;
  isLoading = true;
  hasError = false;
  
  // Добавляем поля для хранения данных, чтобы не пересчитывать их в геттерах
  incomeData: any[] = []; 
  bulletPoints: { negative_factors: string[], positive_factors: string[] } = { negative_factors: [], positive_factors: [] };
  fullName: string = '';
  selectedLimit: number = 0;

  private destroy$ = new Subject<void>();

  constructor(
    private route: ActivatedRoute, 
    private dataService: DataService,
    private cdr: ChangeDetectorRef,
    @Inject(PLATFORM_ID) private platformId: Object
  ) {}

  ngOnInit() {
    if (isPlatformBrowser(this.platformId)) {
      this.route.paramMap.pipe(
        takeUntil(this.destroy$),
        switchMap(params => {
          this.clientID = Number(params.get('id')) ?? null;
          this.isLoading = true;
          this.hasError = false;
          this.scoring = undefined;
          // Очищаем данные графиков при смене клиента
          this.incomeData = [];
          
          console.log('Loading scoring for client:', this.clientID);
          return this.dataService.getScoringClient(this.clientID!);
        })
      ).subscribe({
        next: (scoring) => {
          console.log('Scoring loaded successfully:', scoring);
          this.scoring = scoring;
          
          // !!! ВАЖНО: Подготавливаем данные один раз здесь !!!
          this.prepareViewData();

          this.isLoading = false;
          this.cdr.detectChanges();
        },
        error: (error) => {
          console.error('Error loading scoring:', error);
          this.hasError = true;
          this.isLoading = false;
          this.cdr.detectChanges();
        }
      });
    } else {
      this.isLoading = true;
    }
  }

  ngOnDestroy() {
    this.destroy$.next();
    this.destroy$.complete();
  }

  // Метод для подготовки данных (вызывается один раз при получении ответа)
  private prepareViewData() {
    if (!this.scoring) return;

    // 1. Фиксируем массив для графика
    this.incomeData = [
      { name: 'Текущий доход', value: this.scoring.income || 0 },
      { name: 'Прогнозируемый доход', value: this.scoring.predict_income || 0 },
    ];

    // 2. Фиксируем буллиты
    this.bulletPoints = {
      negative_factors: this.scoring.negative_factors || [],
      positive_factors: this.scoring.positive_factors || []
    };

    // 3. Фиксируем имя
    this.fullName = getFullName({
      first_name: this.scoring.first_name,
      last_name: this.scoring.last_name,
      middle_name: this.scoring.middle_name
    } as Client);

    // 4. Фиксируем лимит
    this.selectedLimit = this.scoring.credit_limit ?? 0;
  }

  get recommendations() {
    return this.scoring?.recommendations;
  }

  protected readonly getFullName = getFullName;
}