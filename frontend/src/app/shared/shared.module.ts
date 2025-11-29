import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { NzButtonModule } from 'ng-zorro-antd/button';
import { NzTableModule } from 'ng-zorro-antd/table';


@NgModule({
  imports: [
    CommonModule,
    NzButtonModule,
    NzTableModule
  ],
  declarations: [],
  exports: [
    NzButtonModule,
    NzTableModule,
  ]
})
export class SharedModule { }
