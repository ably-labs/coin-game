import { ChangeDetectorRef, Component, NgZone, OnInit } from '@angular/core';
import { WalletService } from '../wallet.service';
import { Realtime, Types } from 'ably';
import { MatDialog, MatDialogConfig } from '@angular/material/dialog';
import { BoardModalComponent } from '../board-modal/board-modal.component';
import { environment } from '../../environments/environment';
@Component({
  selector: 'app-home',
  templateUrl: './home.component.html',
  styleUrls: ['./home.component.css']
})
export class HomeComponent implements OnInit {
  client: Realtime = new Realtime({
    key: `${environment.API_KEY}`,
    clientId: "r5ew9P"
  });

  portfolio: any = "";
  bitcoinPrice!: number;
  playerName: string = "";
  player: any = "";
  coinWorth!: number;
  quantity: number = 0;
  enter: boolean = false;
  chanName = "stock";
  channel = this.client.channels.get(this.chanName);
  view = false;
  subscriptions: any = [];
  members: any = [];
  timer: string = "30"
  start: boolean = false;


  constructor(public cd: ChangeDetectorRef, private walletService: WalletService,
    private dialog: MatDialog, public zone: NgZone) {

  }

  ngOnInit() {
    this.getCurrentPrice();

    this.channel.subscribe("buy", (message) => {
      this.processTransaction(message)
    });

    this.channel.subscribe("sell", (message) => {
      this.processTransaction(message)
    });

    this.channel.presence.subscribe("enter", (presence) => {
      this.enterPresence(presence)
    })

    this.channel.presence.subscribe("leave", (presence) => {
      this.leavePresence(presence)
    })

    this.channel.subscribe('start', (message) => this.zone.run(() => {
      this.start = true;
    }))

    this.channel.subscribe('end', (message) => this.zone.run(() => {
      if (this.enter) {
        this.openDialog()
        this.reset()
      }
      this.start = false;
      this.timer = message.data
    }))

    this.channel.subscribe('time', (message) => {
      this.timer = message.data
    })
  }


  getCurrentPrice(): void {
    this.walletService.getBitcoinPrice().subscribe((data) => {
      this.bitcoinPrice = data.data;
      this.coinWorth = +(this.portfolio.CoinQuantity * this.bitcoinPrice).toFixed(2);
      this.cd.detectChanges();
    })
  }

  getWallet(name: string): void {
    this.walletService.getWalletBalance(name).subscribe((data) => {
      this.portfolio = data;
    });
  }

  createPlayer(name: string): void {
    this.enter = true;
    this.walletService.createPlayer(name).subscribe((data) => {
      this.portfolio = data;
    })
    this.channel.presence.enter({ name: name }, (err) => {
    });
  }

  buyCoin(quantity: string, name: string): void {
    this.walletService.buyBitcoin(this.bitcoinPrice, +quantity, name).subscribe((response) => {
      this.getWallet(name)
    });
  }

  sellCoin(quantity: string, name: string): void {
    this.walletService.sellBitcoin(this.bitcoinPrice, +quantity, name).subscribe((response) => {
      this.getWallet(name)
    });

  }

  reset(): void {
    this.channel.presence.leave((err: any) => {
    });
    this.enter = false;
    this.playerName = ""
  }

  openDialog() {
    const dialogConfig = new MatDialogConfig();

    let dialogRef = this.dialog.open(BoardModalComponent, {
      data: {
        result: this.subscriptions.sort((a: any, b: any) => (a.WalletBalance < b.WalletBalance) ? 1 : -1),
        price: this.bitcoinPrice
      }
    });

    dialogRef.afterClosed().subscribe(data => {
      this.reset()
    })
  }

  countdown(time: string) {
    let timeNum = +time
    const interval = setInterval(() => {
      timeNum -= 1
      this.channel.publish("time", `${timeNum}`, (err) => {
        if (err) {
          console.log("error")
        }
    })
      if (+this.timer <= 0) {
        this.channel.publish("end", '30', (err) => {
          if (err) {
            console.log("error")
          }
      })
        clearInterval(interval)
        return
      }

    }, 1000)
  }

  startGame() {
    let data = {
      name: this.playerName,
      time: this.timer
    }
    this.channel.publish("start", data, (err) => {
      if (err) {
        console.log("error")
      }
    })
    this.countdown(this.timer)
  }

  processTransaction(message: Types.Message) {
    let userData = message.data
    const index = this.subscriptions.findIndex(((obj: { Player: any; }) => {
      return obj.Player == userData.Player;
    }));
    if (index === -1) {
      this.subscriptions.push(userData)
    } else {
      this.subscriptions[index] = userData
    }
  }

  enterPresence(presence: Types.PresenceMessage) {
    this.members.push(presence.data)
  }

  leavePresence(presence: Types.PresenceMessage) {
    this.members = this.members.filter((member: { name: any; }) => {
      return member.name != presence.data.name;
    });
    this.subscriptions = this.subscriptions.filter((member: { Player: any; }) => {
      return member.Player != presence.data.name;
    })
  }

}
