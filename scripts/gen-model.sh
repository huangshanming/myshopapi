#!/usr/bin/env bash
# Generate go-zero sqlx models into services/<svc>/internal/modelgen.
# Repositories stay hand-written and consume generated entity structs.
#
# Usage: ./scripts/gen-model.sh [user|catalog|order|merchant|all]
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
export PATH="$PATH:$(go env GOPATH)/bin"
HOME_TPL="${GOCTL_HOME:-$ROOT/deploy/goctl-template}"
TARGET="${1:-all}"

if ! command -v goctl >/dev/null 2>&1; then
  echo "goctl not found; install: go install github.com/zeromicro/go-zero/tools/goctl@v1.8.5" >&2
  exit 1
fi

postprocess_modelgen() {
  python3 - "$1" <<'PY'
import re, sys
from pathlib import Path
dir = Path(sys.argv[1])
RENAMES = {
    'Users': 'User', 'UserAddresses': 'UserAddress', 'UserWallets': 'UserWallet', 'UserWalletLogs': 'UserWalletLog',
    'UserPointLogs': 'UserPointLog', 'TaskDefinitions': 'TaskDefinition', 'PointsProducts': 'PointsProduct',
    'PointsExchangeOrders': 'PointsExchangeOrder', 'UserNotifications': 'UserNotification',
    'UserNotificationBatches': 'UserNotificationBatch', 'Regions': 'Region',
    'ProductCategories': 'ProductCategory', 'Products': 'Product', 'ProductSkus': 'ProductSku',
    'ProductImages': 'ProductImage', 'ProductTags': 'ProductTag', 'ProductTagRels': 'ProductTagRel',
    'ProductAttrTemplates': 'ProductAttrTemplate', 'ProductAttrs': 'ProductAttr',
    'ProductSchedules': 'ProductSchedule', 'ProductBatchJobs': 'ProductBatchJob', 'ProductOpLogs': 'ProductOpLog',
    'ProductFavorites': 'ProductFavorite', 'ShopRoles': 'ShopRole', 'ShopMenus': 'ShopMenu',
    'ShopRoleMenus': 'ShopRoleMenu', 'ShopUserRoles': 'ShopUserRole',
    'CommunityCommentEmojis': 'CommunityCommentEmoji', 'ArticleLikes': 'ArticleLike',
    'ArticleFavorites': 'ArticleFavorite', 'ArticleAudiences': 'ArticleAudience',
    'HomepageBanners': 'HomepageBanner', 'ShopNotifications': 'ShopNotification',
    'Orders': 'Order', 'OrderItems': 'OrderItem', 'OrderAfterSales': 'OrderAfterSale',
    'ProductReviews': 'ProductReview', 'ProductReviewImages': 'ProductReviewImage',
    'LogisticsCompanies': 'LogisticsCompany', 'Shops': 'Shop', 'ShopApplications': 'ShopApplication',
    'ShopMembers': 'ShopMember', 'ShopWallets': 'ShopWallet', 'ShopWalletLogs': 'ShopWalletLog',
    'SeckillRules': 'SeckillRule', 'SeckillSessions': 'SeckillSession', 'SeckillEntries': 'SeckillEntry',
    'HomepageSlotPackages': 'HomepageSlotPackage', 'HomepageSlotSettings': 'HomepageSlotSetting',
    'HomepageSlotOrders': 'HomepageSlotOrder', 'HomepageThemeSlots': 'HomepageThemeSlot',
    'HomepageThemePackages': 'HomepageThemePackage', 'HomepageThemeOrders': 'HomepageThemeOrder',
    'Coupons': 'Coupon', 'CouponScopes': 'CouponScope', 'UserCoupons': 'UserCoupon',
    'CouponGrants': 'CouponGrant', 'CouponRedeemLogs': 'CouponRedeemLog',
}
def fix_initialisms(name: str) -> str:
    for a, b in [('Id','ID'),('Url','URL'),('Uri','URI'),('Sku','SKU'),('Ip','IP'),('Api','API'),('Uuid','UUID'),('Json','JSON'),('Xml','XML')]:
        name = re.sub(rf'({a})(?=[A-Z]|$)', b, name)
    return name
for f in dir.rglob('*.go'):
    text = f.read_text(); orig = text
    if 'time.Time' in text:
        text = text.replace('time.Time', 'common.LocalTime')
        if re.search(r'\btime\.', text) is None:
            text = re.sub(r'\n\t"time"\n', '\n', text)
        if 'mymall/common' not in text and 'common.LocalTime' in text:
            text = text.replace('import (\n', 'import (\n\t"mymall/common"\n', 1)
    text = re.sub(r'(\s+)sql\.NullString(\s+`db:)', r'\1string\2', text)
    text = re.sub(r'(\s+)sql\.NullInt64(\s+`db:)', r'\1uint64\2', text)
    text = text.replace('sql.NullTime', 'common.LocalTime')
    text = re.sub(r'(\s+)int64(\s+`db:)', r'\1int\2', text)
    for old, new in sorted(RENAMES.items(), key=lambda kv: -len(kv[0])):
        if old != new:
            text = re.sub(rf'\b{old}\b', new, text)
    def field_line(m):
        return m.group(1) + fix_initialisms(m.group(2)) + m.group(3)
    text = re.sub(r'(^\s+)([A-Z][A-Za-z0-9]*)(\s+\S+\s+`db:)', field_line, text, flags=re.M)
    # Keep the dot: data.UserId -> data.UserID
    text = re.sub(r'\.([A-Z][A-Za-z0-9]*)', lambda m: '.' + fix_initialisms(m.group(1)), text)
    text = re.sub(r'\bdata([A-Z][A-Za-z0-9]*)\b', r'data.\1', text)
    text = re.sub(r'\bnewData([A-Z][A-Za-z0-9]*)\b', r'newData.\1', text)
    # Add snake_case json tags from db column names (password omitted from JSON).
    def add_json(m):
        tag = m.group(1)
        if 'json:' in tag:
            return '`' + tag + '`'
        dm = re.search(r'db:"([^"]+)"', tag)
        if not dm or dm.group(1) == '-':
            return '`' + tag + '`'
        col = dm.group(1)
        jn = '-' if col == 'password' else col
        return '`' + tag + f' json:"{jn}"`'
    text = re.sub(r'`([^`]*db:"[^"]+"[^`]*)`', add_json, text)
    if 'sql.' not in text and '"database/sql"' in text:
        text = re.sub(r'\n\t"database/sql"\n', '\n', text)
    if text != orig:
        f.write_text(text)
PY
}

write_composites() {
  local svc="$1"
  local dir="$ROOT/services/${svc}-service/internal/modelgen"
  case "$svc" in
    user)
      cat >"$dir/sys_role_menu.go" <<'EOF'
package modelgen

type SysRoleMenu struct {
	RoleID uint64 `db:"role_id"`
	MenuID uint64 `db:"menu_id"`
}

func (SysRoleMenu) TableName() string { return "sys_role_menu" }
EOF
      cat >"$dir/sys_user_role.go" <<'EOF'
package modelgen

type SysUserRole struct {
	UserID uint64 `db:"user_id"`
	RoleID uint64 `db:"role_id"`
}

func (SysUserRole) TableName() string { return "sys_user_role" }
EOF
      ;;
    catalog)
      cat >"$dir/product_tag_rels.go" <<'EOF'
package modelgen

type ProductTagRel struct {
	ProductID uint64 `db:"product_id"`
	TagID     uint64 `db:"tag_id"`
}

func (ProductTagRel) TableName() string { return "product_tag_rels" }
EOF
      cat >"$dir/shop_role_menus.go" <<'EOF'
package modelgen

type ShopRoleMenu struct {
	RoleID uint64 `db:"role_id"`
	MenuID uint64 `db:"menu_id"`
}

func (ShopRoleMenu) TableName() string { return "shop_role_menus" }
EOF
      cat >"$dir/shop_user_roles.go" <<'EOF'
package modelgen

type ShopUserRole struct {
	ShopID uint64 `db:"shop_id"`
	UserID uint64 `db:"user_id"`
	RoleID uint64 `db:"role_id"`
}

func (ShopUserRole) TableName() string { return "shop_user_roles" }
EOF
      ;;
  esac
}

gen_one() {
  local name="$1"
  local ddl="$ROOT/scripts/ddl/${name}-tables.sql"
  local out="$ROOT/services/${name}-service/internal/modelgen"
  [[ -f "$ddl" ]] || { echo "missing $ddl" >&2; exit 1; }
  echo "==> goctl model: $name"
  rm -rf "$out"
  mkdir -p "$out"
  python3 - "$ddl" "$out" "$HOME_TPL" <<'PY'
import re, subprocess, sys
from pathlib import Path
ddl, out, home = Path(sys.argv[1]), Path(sys.argv[2]), sys.argv[3]
text = ddl.read_text()
ok = fail = skip = 0
for m in re.finditer(r'CREATE TABLE IF NOT EXISTS (\w+)\s*\((.*?)\)\s*ENGINE\s*=\s*InnoDB[^;]*;', text, re.S|re.I):
    name, body = m.group(1), m.group(2)
    pk = re.search(r'PRIMARY KEY\s*\(([^)]+)\)', body, re.I)
    if pk and len([c for c in pk.group(1).split(',') if c.strip()]) > 1:
        print(f'    skip composite {name}')
        skip += 1
        continue
    one = Path(f'.tmp/one-ddl/{name}.sql')
    one.parent.mkdir(parents=True, exist_ok=True)
    one.write_text(m.group(0) + '\n')
    r = subprocess.run(
        ['goctl', 'model', 'mysql', 'ddl', '-src', str(one), '-dir', str(out),
         '-style', 'go_zero', '--home', home],
        capture_output=True, text=True)
    if r.returncode != 0:
        print(f'    FAIL {name}: {(r.stderr or r.stdout)[-200:]}')
        fail += 1
    else:
        ok += 1
print(f'    ok={ok} fail={fail} skip={skip}')
if fail:
    raise SystemExit(1)
PY
  write_composites "$name"
  postprocess_modelgen "$out"
  echo "    -> $out"
}

case "$TARGET" in
  user|catalog|order|merchant) gen_one "$TARGET" ;;
  all)
    gen_one user
    gen_one catalog
    gen_one order
    gen_one merchant
    ;;
  *)
    echo "usage: $0 [user|catalog|order|merchant|all]" >&2
    exit 2
    ;;
esac

echo "goctl model done."
