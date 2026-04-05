// Polyfill: newer protoc-gen-js generates readStringRequireUtf8() calls, but the
// google-protobuf npm runtime only ships readString(). They are equivalent for our use.
import * as jspb from 'google-protobuf';
if (!(jspb.BinaryReader.prototype as any).readStringRequireUtf8) {
  (jspb.BinaryReader.prototype as any).readStringRequireUtf8 = jspb.BinaryReader.prototype.readString;
}

import { platformBrowserDynamic } from '@angular/platform-browser-dynamic';

import { AppModule } from './app/app.module';

platformBrowserDynamic()
  .bootstrapModule(AppModule)
  .catch((err) => console.error(err));
