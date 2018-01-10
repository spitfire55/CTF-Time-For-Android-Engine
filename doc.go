// Copyright 2017-2018 Dale Lakes <spitfire@spitfy.re>. All rights reserved.
// Use of this source code is governed by the MIT license located in the LICENSE file.

/*
Package goctftime parses and stores data from ctftime.org into a Firestore database so that it can be easily indexed and
queried by an Android application that displays CTF Time data.

Handlers

The various handlers listen on the server for a GET request to their respective path. This triggers the execution of the
handler logic. Each handler is responsible for parsing and storing some portion of ctftime.org. Handler logic is broken up
into two phases.

The first phase triggers multiple goroutines to parse and store data concurrently. By default, the maximum number
of goroutines running at once is 10. To change the maximum number of goroutines running at once, modify the maxRoutines
variable. This concurrent phase only requests pages that we have scraped before.

The second phase operates on a single thread and checks to see if new content exists. If new content exists, it is
parsed and stored in Firestore. Finally, we update the value used in phase one to delineate the range of known content.

GetLast and UpdateLast

Each scraper operates concurrently as it requests known content. The definition of "known" is based on a Firestore value for
the respective piece of content (teams, events, etc.). This value is retrieved and possibly updated on each scraping
iteration.
*/
package goctftime
