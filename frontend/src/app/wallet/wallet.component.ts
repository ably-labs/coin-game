import { ChangeDetectorRef, Component, OnInit, Type } from '@angular/core';
import { WalletService } from '../wallet.service';

@Component({
  selector: 'app-wallet',
  templateUrl: './wallet.component.html',
  styleUrls: ['./wallet.component.css']
})
export class WalletComponent implements OnInit {

  walletBalance!: number;
  bitcoinPrice: any = "";

  constructor(public cd: ChangeDetectorRef, private walletService: WalletService) {

   }

  ngOnInit() {
    this.getWallet();
    this.getCurrentPrice();
  }


  getCurrentPrice(): void {
    this.walletService.getBitcoinPrice().subscribe((data) => {
      this.bitcoinPrice = data.data;
      this.cd.detectChanges();
    })
  }

  getWallet(): void {
    this.walletService.getWalletBalance().subscribe((data) => {
        this.walletBalance = +data;
    });
  }

  buyCoin(): void {
    let quantity = 1;
    this.walletService.buyBitcoin(this.bitcoinPrice, quantity).subscribe((response) => {
      this.getWallet()
      console.log(typeof(response), "here it goes")
    });
  }

  sellCoin(): void {
    let quantity = 1;
    this.walletService.buyBitcoin(this.bitcoinPrice, quantity).subscribe((response) => {
      this.getWallet()
      console.log(response, "here it goes")
    });
  }

}
