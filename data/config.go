package data

import (
	"appengine"
	"appengine/datastore"
	"log"
)

/*
 * K/V values:
 *  lastCheck Time (as string)
 *  bookmarkUrl string
 *  bookmarkService string (default to delicious)
 *  headings string (comma delimited)
 */
type KVPair struct {
	Key string;
	Value string;
}

func configOptionKey (c appengine.Context) *datastore.Key {
	return datastore.NewKey(c, "KVPair", "config_options", 0, nil)
}

func GetConfigOption (c appengine.Context, key string) (value string) {
	q := datastore.NewQuery("KVPair").
		Filter("Key = ", key);
	results := q.Run(c)
	for {
		var kv KVPair
		_, err := results.Next(&kv);
		if err == datastore.Done {
			break // we're done
		}
		if err != nil {
			log.Fatal("error fetching values: ", err)
			break
		}
		value = kv.Value
	}
	return value
}

func SetConfigOption (c appengine.Context, key string, value string) {
	q := datastore.NewQuery("KVPair").
		Filter("Key = ", key);
	results := q.Run(c)

	db_entry := datastore.NewIncompleteKey(c, "KVPair", configOptionKey(c));

	var kv KVPair
	for {
		db_key, err := results.Next(&kv);
		if err == datastore.Done {
			break // the result was not found
		}
		if err != nil {
			log.Fatal("error fetching values: ", err)
			break
		}
		db_entry = db_key
	}
	kv.Key   = key;
	kv.Value = value;

	_, err := datastore.Put(c, db_entry, &kv);
	if err != nil {
		log.Fatal("error putting option: ", err);
	}
}

func GetAllConfigOptions (c appengine.Context) (options []KVPair) {
	q := datastore.NewQuery("KVPair").
		Order("Key");

	_, err := q.GetAll(c, &options);

	if err != nil {
		log.Fatal("error fetching options: ", err);
	}

	return options
}

func DeleteConfigOption (c appengine.Context, key string) {
	q := datastore.NewQuery("KVPair").
		Filter("Key = ", key);
	results := q.Run(c)

	var kv KVPair
	for {
		db_key, err := results.Next(&kv);
		if err == datastore.Done {
			break // the result was not found
		}
		if err != nil {
			log.Fatal("error fetching values: ", err)
			break
		}
		// the option was found, delete it.
		err = datastore.Delete(c, db_key);
		if err != nil {
			log.Fatal("error deleting key: ", err)
		}
	}
}
