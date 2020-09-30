package km

const(
	initDzzz=`create table if not exists dzzz(
		bh	text,	-- 定制转账编号
		xh	text,	-- 定制转账序号
		mc	text,	-- 名称
		czjg	text,	-- 操作机构号
		czjglx	text,	-- 操作机构类型
		czjgfh	text,	-- 操作机构所在分行
		czlwjg	text,	-- 操作机构例外机构
		bz		text,	-- 币种
		jdbz	text,	-- 借贷标志
		zhjg	text,	-- 账户所在机构码
		zhjglx	text,	-- 账户机构类型
		zhjgfh	text,	-- 账户所在分行
		zhlwjg	text,	-- 账户例外机构
		km		text,	-- 科目
		xh		int,	-- 序号
		yxkjg	text,	-- 是否允许跨机构
		yxhz	text,	-- 是否允许红字
		primary key(bh,xh)
`
)