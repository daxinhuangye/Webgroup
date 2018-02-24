/**
 * Created by Administrator on 2018/2/24.
 */
app.filter('Tohtml', ['$sce', function ($sce) {
    return function (text) {
        text = text || "";
        console.log(text);
       // text = text.replace(/<img src=\/?.*?>/g, '[图片名]', $1);
        return $sce.trustAsHtml(text);
    };
}]);