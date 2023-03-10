#!/bash/bin
# 同步生成 entgo 模型 代码
# 运行请在项目根目录

# 配置数据库连接信息
dsn="mysql://root:Lsxd123.@tcp(192.168.6.71:3308)/cmbi"

SHELL_FOLDER=$(cd "$(dirname "$0")";pwd)
echo $SHELL_FOLDER

entPath="${SHELL_FOLDER}/internal/data/ent"

if [ "${1}" != "" ]; then
  if [ -f "${entPath}/schema/${1}.go" ]; then
    echo "文件已存在，如果要重新生成，请先删除文件：rm -f internal/data/ent/schema/${1}.go"
    exit
  fi

  # 注意配置数据库连接信息、root:Lsxd123.@tcp(192.168.6.71:3308)/cmbi
  entimport -dsn ${dsn} -schema-path="${entPath}/schema" -tables=$1
  #entimport 生成出来不知道为啥有 gorm.io 的东西
  sed -i '' 's/gorm.io\/gen\/field/entgo.io\/ent\/schema\/field/' $entPath/schema/$1.go
  #将所有的 整型 统一为 Int
  sed -i '' 's/field.Uint/field.Int/;s/field.Uint8/field.Int/;s/field.Uint16/field.Int/;s/field.Uint32/field.Int/;s/field.Uint64/field.Int/;s/field.Int8/field.Int/;s/field.Int16/field.Int/;s/field.Int32/field.Int/;s/field.Int64/field.Int/' $entPath/schema/$1.go
  echo "${1} schema 生成完成!!!"
fi

echo "正根据schema 文件生成模型代码..."
go generate ${entPath}
echo "代码模型生成完成!"