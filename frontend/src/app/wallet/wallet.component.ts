import { ChangeDetectionStrategy, Component, Input, OnInit } from '@angular/core';
import { WalletService } from '../wallet.service';
import { Realtime } from 'ably';
import { Observable, retry, Subject, takeUntil } from 'rxjs';

@Component({
  selector: 'app-wallet',
  templateUrl: './wallet.component.html',
  styleUrls: ['./wallet.component.css']
})
export class WalletComponent implements OnInit {

  walletBalance!: number;
  bitcoinPrice: any = "";
  destroyed$ = new Subject();

  constructor(private walletService: WalletService) {

   }

  ngOnInit() {
    this.getWallet();

    // this.walletService.getCurrentPrice().pipe(
    //   retry(),
    //   takeUntil(this.destroyed$)
    // ).subscribe((data) => {
    //   console.log(data, "message above")
    //   this.bitcoinPrice = data;
    // }
    // ,(error) => console.error(error),
    // () => console.log('completed'));

    this.walletService.channel$.subscribe((data) => {
      console.log(data.data, "freaking data lls");
      this.bitcoinPrice = "data.data";
    })

    // this.bitcoinPrice = this.walletService.getCurrentPrice();
    // console.log(this.bitcoinPrice, "line 22")
  }

  ngOnDestroy():void {
    this.destroyed$.next('toh');
  }

  // ngAfterContentInit() {
  //   this.walletService.channel$.subscribe((data) => {
  //     console.log(data.data, "freaking data lls");
  //     this.bitcoinPrice = data.data;
  //   })

  // }

  // getCurrentPrice() {
  //   let ably = new Realtime('VbpYdQ.79UpAA:-lnejxoRLhS_hDPgNrE5XqweLrsLdH0vMZwSQtaKlLI');
  //   let chanName = '[product:ably-coindesk/bitcoin]bitcoin:usd';
  //   let channel = ably.channels.get(chanName);

  //   channel.subscribe((message: any): void => {
  //     this.bitcoinPrice = message.data;
  //     console.log('Received msg', this.bitcoinPrice);
  //   });
  // }

  getWallet(): void {
    this.walletService.getWalletBalance().subscribe((data) => {
        this.walletBalance = +data;
        console.log('tog msg', this.bitcoinPrice);
    });
  }

}
