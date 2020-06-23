package jym

var initJymSQL = `
CREATE TABLE IF NOT EXISTS jym(
    jym     text primary key,   --交易码
    jymc    text,   --交易名称
    jyz     text,   --交易组
    yxj     text,   --优先级
    wdsqjb  text,   --网点授权级别
    zssqjb  text,   --中心授权级别
    wdsq    text,   --网点授权
    zssqjg  text,   --中心授权机构
    zssq    text,   --中心授权
    jnjb    text,   --技能级别
    xzbz    text,   --现转标志
    wb      text,   --外包
    dets    text,   --大额提示
    dzdk    text,   --电子底卡
    sxf     text,   --手续费
    htjc    text,   --后台监测
    szjd    text,   --事中监督
    bssx    text,   --补扫时限
    sc      text,   --审查
    mz      text,   --抹账
    cesq    text,   --超额授权
    fjjyz   text    --辅加交易组
    --shbs    text default "TRUE",   --事后补扫
    --cdjy    text default "FALSE"   --磁道校验
);`

var loadJymSQL = ``
