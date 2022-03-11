import { Component, Inject, OnInit } from '@angular/core';
import { MatDialogRef, MAT_DIALOG_DATA} from '@angular/material/dialog';

@Component({
  selector: 'app-board-modal',
  templateUrl: './board-modal.component.html',
  styleUrls: ['./board-modal.component.css']
})



export class BoardModalComponent implements OnInit {

  constructor(
    public dialogRef: MatDialogRef<BoardModalComponent>,
    @Inject(MAT_DIALOG_DATA) public data: any)
  {
    dialogRef.disableClose = true;
  }
  ngOnInit() {
  }

  restart() {
    this.dialogRef.close()
  }
}
