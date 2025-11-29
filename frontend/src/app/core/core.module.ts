import { NgModule, Optional, SkipSelf } from '@angular/core';
import { CommonModule } from '@angular/common';
import { DataService } from "./services/data.service";
import { ApiService } from './services/api.service';

@NgModule({
  imports: [CommonModule],
  providers: [
    DataService,
    ApiService,
  ]
})
export class CoreModule {
  constructor(@Optional() @SkipSelf() parentModule: CoreModule) {
    if (parentModule) {
      throw new Error('CoreModule уже загружен. Импортируйте только в AppModule');
    }
  }
}
