# upsert and updated_at



### Automatic CreatedAt/UpdatedAt

If your generated SQLBoiler models package can find columns with the
names `created_at` or `updated_at` it will automatically set them
to `time.Now()` in your database, and update your object appropriately.
To disable this feature use `--no-auto-timestamps`.

Note: You can set the timezone for this feature by calling `boil.SetLocation()`

#### Customizing the timestamp columns

Set the `auto-columns` map in your configuration file

```toml
[auto-columns]
    created = "createdAt"
    updated = "updatedAt"
```

#### Skipping Automatic Timestamps

If for a given query you do not want timestamp columns to be re-computed prior
to an insert or update then you can use `boil.SkipTimestamps` on the context you
pass in to the query to prevent them from being updated.

Keep in mind this has no effect on whether or not the column is included in the
insert/update, it simply stops them from being set to `time.Now()` in the struct
before being sent to the database (if they were going to be sent).

#### Overriding Automatic Timestamps

* **Insert**
  * Timestamps for both `updated_at` and `created_at` that are zero values will be set automatically.
  * To set the timestamp to null, set `Valid` to false and `Time` to a non-zero value.
  This is somewhat of a work around until we can devise a better solution in a later version.
* **Update**
  * The `updated_at` column will always be set to `time.Now()`. If you need to override
  this value you will need to fall back to another method in the meantime: `queries.Raw()`,
  overriding `updated_at` in all of your objects using a hook, or create your own wrapper.
* **Upsert**
  * `created_at` will be set automatically if it is a zero value, otherwise your supplied value
  will be used. To set `created_at` to `null`, set `Valid` to false and `Time` to a non-zero value.
  * The `updated_at` column will always be set to `time.Now()`.

### Automatic DeletedAt (Soft Delete)



