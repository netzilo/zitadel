import { Component, OnInit, ViewChild } from '@angular/core';
import { MatSort, Sort } from '@angular/material/sort';
import { MatLegacyTableDataSource as MatTableDataSource } from '@angular/material/legacy-table';
import { BehaviorSubject, Observable } from 'rxjs';
import { ListEventsRequest, ListEventsResponse } from 'src/app/proto/generated/zitadel/admin_pb';
import { AdminService } from 'src/app/services/admin.service';
import { Breadcrumb, BreadcrumbService, BreadcrumbType } from 'src/app/services/breadcrumb.service';
import { Event } from 'src/app/proto/generated/zitadel/event_pb';
import { PaginatorComponent } from 'src/app/modules/paginator/paginator.component';
import { LiveAnnouncer } from '@angular/cdk/a11y';
import { ToastService } from 'src/app/services/toast.service';
import { ConnectedPosition, ConnectionPositionPair } from '@angular/cdk/overlay';
import { ActionKeysType } from 'src/app/modules/action-keys/action-keys.component';
import { DisplayJsonDialogComponent } from 'src/app/modules/display-json-dialog/display-json-dialog.component';
import { MatLegacyDialog as MatDialog } from '@angular/material/legacy-dialog';

enum EventFieldName {
  EDITOR = 'editor',
  AGGREGATE = 'aggregate',
  SEQUENCE = 'sequence',
  CREATIONDATE = 'creationDate',
  TYPE = 'type',
  PAYLOAD = 'payload',
}

@Component({
  selector: 'cnsl-events',
  templateUrl: './events.component.html',
  styleUrls: ['./events.component.scss'],
})
export class EventsComponent implements OnInit {
  public INITPAGESIZE = 20;
  public sortAsc = false;

  public displayedColumns: string[] = [
    EventFieldName.TYPE,
    EventFieldName.AGGREGATE,
    EventFieldName.EDITOR,
    EventFieldName.SEQUENCE,
    EventFieldName.CREATIONDATE,
    EventFieldName.PAYLOAD,
  ];

  public currentRequest: ListEventsRequest = new ListEventsRequest();

  @ViewChild(MatSort) public sort!: MatSort;
  @ViewChild(PaginatorComponent) public paginator!: PaginatorComponent;
  public dataSource: MatTableDataSource<Event> = new MatTableDataSource<Event>([]);

  public _done: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(false);
  public done: Observable<boolean> = this._done.asObservable();

  public _loading: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(false);

  private _data: BehaviorSubject<Event[]> = new BehaviorSubject<Event[]>([]);

  constructor(
    private adminService: AdminService,
    private breadcrumbService: BreadcrumbService,
    private _liveAnnouncer: LiveAnnouncer,
    private toast: ToastService,
    private dialog: MatDialog,
  ) {
    const breadcrumbs = [
      new Breadcrumb({
        type: BreadcrumbType.INSTANCE,
        name: 'Instance',
        routerLink: ['/instance'],
      }),
    ];
    this.breadcrumbService.setBreadcrumb(breadcrumbs);
  }

  ngOnInit(): void {
    const req = new ListEventsRequest();
    req.setLimit(this.INITPAGESIZE);

    this.loadEvents(req);
  }

  public loadEvents(filteredRequest: ListEventsRequest): Promise<void> {
    this._loading.next(true);

    this.currentRequest = filteredRequest;
    console.log('load', this.currentRequest.toObject());

    return (
      this.adminService
        .listEvents(this.currentRequest)
        //   .then((resp) => {
        //     const eventsList = resp.getEventsList();
        //     eventsList.map((event) => {
        //       const payload = event.getPayload();
        //       const pl = payload?.toJavaScript();
        //       console.log(pl);
        //     });
        //     return resp.toObject();
        //   })
        .then((res: ListEventsResponse) => {
          const eventList = res.getEventsList();
          this._data.next(eventList);

          const concat = this.dataSource.data.concat(eventList);
          this.dataSource = new MatTableDataSource<Event>(concat);

          this._loading.next(false);

          if (eventList.length === 0) {
            this._done.next(true);
          } else {
            this._done.next(false);
          }
        })
        .catch((error) => {
          console.error(error);
          this._loading.next(false);
          this._data.next([]);
        })
    );
  }

  public refresh(): void {
    const req = new ListEventsRequest();
    req.setLimit(this.paginator.pageSize);
  }

  public sortChange(sortState: Sort) {
    if (sortState.direction && sortState.active) {
      this.dataSource = new MatTableDataSource<Event>([]);

      this._liveAnnouncer.announce(`Sorted ${sortState.direction}ending`);
      this.sortAsc = sortState.direction === 'asc';

      this.currentRequest = new ListEventsRequest();
      this.currentRequest.setLimit(this.INITPAGESIZE);
      this.currentRequest.setAsc(this.sortAsc ? true : false);

      this.loadEvents(this.currentRequest);
    } else {
      this._liveAnnouncer.announce('Sorting cleared');
    }
  }

  public openDialog(event: Event): void {
    this.dialog.open(DisplayJsonDialogComponent, {
      data: {
        event: event,
      },
      width: '450px',
    });
  }

  public more(): void {
    const sequence = this.getCursor();
    this.currentRequest.setSequence(sequence);
    this.loadEvents(this.currentRequest);
  }

  public filterChanged(filterRequest: ListEventsRequest) {
    this.dataSource = new MatTableDataSource<Event>([]);

    this.currentRequest = new ListEventsRequest();
    this.currentRequest.setLimit(this.INITPAGESIZE);
    this.currentRequest.setAsc(this.sortAsc ? true : false);

    this.currentRequest.setAggregateTypesList(filterRequest.getAggregateTypesList());
    this.currentRequest.setAggregateId(filterRequest.getAggregateId());
    this.currentRequest.setEventTypesList(filterRequest.getEventTypesList());
    this.currentRequest.setEditorUserId(filterRequest.getEditorUserId());

    this.loadEvents(this.currentRequest);
  }

  private getCursor(): number {
    const current = this._data.value;

    if (current.length) {
      const sequence = current[current.length - 1].toObject().sequence;
      return sequence;
    }
    return 0;
  }
}
