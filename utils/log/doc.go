// Copyright 2015, Cyrill @ Schumacher.fm and the CoreStore contributors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

/*
Package log contains NullLogger, StdLogger and the Logger interface.

Logging

Interface Logger is used all over the place and there are no other dependencies.
Default Logger is a null logger. You must take care to implement a logger which
is also thread safe.

Convention: Because recording a human-meaningful message is common and good
practice, the first argument to every logging method is the value to the
*implicit* key 'msg'. You may supply any additional context as a set of
key/value pairs to the logging function.

To initialize your own logger you must somewhere set the logging object to the
util/log package.

	import "github.com/corestoreio/csfw/utils/log"

	func init() {
		log.Set(NewMyCustomLogger())
	}

Level guards exists to avoid the cost of building arguments. Get in the
habit of using guards.

	import "github.com/corestoreio/csfw/utils/log"

	if log.IsDebug() {
		log.Debug("some message", "key1", expensive())
	}

Standardizes on key-value pair argument sequence:

	import "github.com/corestoreio/csfw/utils/log"

	log.Debug("message from inside Fn()", "key1", value1, "key2", value2)

	// instead of this
	log.WithFields(logrus.Fields{"m": "pkg", "key1": value1, "key2": value2}).Debug("inside fn()")

Please consider the key-value pairs when implementing your own logger.

Recommended Loggers are https://github.com/mgutz/logxi and https://github.com/Sirupsen/logrus
and https://github.com/inconshreveable/log15

Standard Logger

CoreStore provides a leveled logger based on Go's standard library without any
dependencies. This StdLogger obeys to the interface Logger of this package.

	import "github.com/corestoreio/csfw/utils/log"

	func init() {
		log.Set(log.NewStdLogger())
	}

log.NewStdLogger() accepts a wide range of optional arguments. Please see the functions Std*Option().

*/
package log
