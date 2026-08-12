# MySQL 初始化（问数相关）

仓库主库仍使用宿主机 MySQL（homestead :3306）。
可在此目录放置 askdata 专用 DDL / 种子 SQL，例如：

- `01_schema.sql`
- `02_seed.sql`

本地导入示例：

```bash
mysql -uhomestead -psecret mymall < docker/mysql/01_schema.sql
```
