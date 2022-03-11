import { HttpClientModule, HTTP_INTERCEPTORS } from '@angular/common/http';
import { NgModule } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';
import { FormsModule } from '@angular/forms';

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { HomeComponent } from './home/home.component';
import { WalletService } from './wallet.service';
import { WithCredentialsInterceptor } from './interceptor';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { MatButtonModule} from '@angular/material/button'
import { MatInputModule} from '@angular/material/input';
import { MatTableModule} from '@angular/material/table';
import {MatListModule} from '@angular/material/list';
import { MatCardModule } from '@angular/material/card';
import { MatDialogModule } from '@angular/material/dialog';
import { BoardModalComponent } from './board-modal/board-modal.component';
import { TimerPipe } from './timer.pipe';

@NgModule({
  declarations: [
    AppComponent,
    HomeComponent,
    BoardModalComponent,
    TimerPipe,
  ],
  imports: [
    BrowserModule,
    AppRoutingModule,
    HttpClientModule,
    FormsModule,
    BrowserAnimationsModule,
    MatButtonModule,
    MatInputModule,
    MatTableModule,
    MatListModule,
    MatCardModule,
    MatDialogModule
  ],
  providers: [WalletService,
  {
    provide: HTTP_INTERCEPTORS,
    useClass: WithCredentialsInterceptor,
    multi: true
  }],
  bootstrap: [AppComponent]
})
export class AppModule { }
