class StatementQueryParamsDto {
  final int? limit;
  final String? cursor;
  final String? cursorId;
  final String? from;
  final String? to;

  const StatementQueryParamsDto({
    this.limit,
    this.cursor,
    this.cursorId,
    this.from,
    this.to,
  });

  Map<String, dynamic> toMap() {
    final map = <String, dynamic>{};
    if (limit != null) map['limit'] = limit;
    if (cursor != null) map['cursor'] = cursor;
    if (cursorId != null) map['cursor_id'] = cursorId;
    if (from != null) map['from'] = from;
    if (to != null) map['to'] = to;
    return map;
  }
}
