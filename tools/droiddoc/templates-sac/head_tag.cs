<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<?cs if:page.metaDescription ?>
<meta name="Description" content="<?cs var:page.metaDescription ?>">
<?cs /if ?>
<link rel="shortcut icon" type="image/x-icon" href="<?cs var:toroot ?>favicon.ico" />
<title><?cs 
  if:page.title ?><?cs 
    var:page.title ?> | <?cs
  /if ?>Android Developers</title>

<!-- STYLESHEETS -->
<link rel="stylesheet"
href="<?cs if:android.whichdoc != 'online' ?>http:<?cs /if ?>//fonts.googleapis.com/css?family=Roboto:regular,medium,thin,italic,mediumitalic,bold" title="roboto">
<link href="<?cs var:toroot ?>assets/css/default.css" rel="stylesheet" type="text/css">

<?cs if:reference && !(reference.gms || reference.gcm || sac) ?>
<!-- FULLSCREEN STYLESHEET -->
<link href="<?cs var:toroot ?>assets/css/fullscreen.css" rel="stylesheet" class="fullscreen"
type="text/css">
<?cs /if ?>

<!-- JAVASCRIPT -->
<script src="<?cs if:android.whichdoc != 'online' ?>http:<?cs /if ?>//www.google.com/jsapi" type="text/javascript"></script>
<?cs
if:devsite
  ?><script src="<?cs var:toroot ?>_static/js/android_3p-bundle.js" type="text/javascript"></script><?cs
else
  ?><script src="<?cs var:toroot ?>assets/js/android_3p-bundle.js" type="text/javascript"></script><?cs
/if ?>
<script type="text/javascript">
  var toRoot = "<?cs var:toroot ?>";
  <?cs if:devsite ?>
  var devsite = true;
  <?cs else ?>
  var devsite = false;
  <?cs /if ?>
</script>
<script src="<?cs var:toroot ?>assets/js/docs.js" type="text/javascript"></script>
<?cs if:reference.gms || reference.gcm || google?>
<script src="<?cs var:toroot ?>gms_navtree_data.js" async type="text/javascript"></script>
<script src="<?cs var:toroot ?>gcm_navtree_data.js" async type="text/javascript"></script>
<?cs elif:devices ?>
<script src="<?cs var:toroot ?>navtree_data.js" type="text/javascript"></script>
<?cs elif:reference ?>
<script src="<?cs var:toroot ?>navtree_data.js" async type="text/javascript"></script>
<?cs /if ?>

<script type="text/javascript">
  var _gaq = _gaq || [];
  _gaq.push(['_setAccount', 'UA-45455297-1']);
  _gaq.push(['_trackPageview']);

  (function() {
    var ga = document.createElement('script'); ga.type = 'text/javascript'; ga.async = true;
    ga.src = ('https:' == document.location.protocol ? 'https://ssl' : 'http://www') + '.google-analytics.com/ga.js';
    var s = document.getElementsByTagName('script')[0]; s.parentNode.insertBefore(ga, s);
  })();
</script>

    <script type="text/javascript" src="https://www.google.com/jsapi"></script>
    <script type="text/javascript">
      google.load("visualization", "1", {packages:["corechart"]});
      google.setOnLoadCallback(drawChart);
      function drawChart() {
        var data = google.visualization.arrayToDataTable([
          ['Project', 'Number of commits'],
          ['root', 10866126], ['external', 3824855], ['prebuilts', 1978812],
          ['frameworks', 1258788], ['ndk', 1213050], ['packages', 1026796],
          ['cts', 479854], ['hardware', 390558], ['libcore', 250101],
          ['device', 110041], ['development', 92450], ['system', 68885],
          ['bionic', 64332], ['developers', 45389], ['sdk', 23251],
          ['dalvik', 16062], ['build', 8122], ['pdk', 7149],
          ['docs', 3493], ['bootable', 2150], ['libnativehelper', 1938],
                  ['abi', 28], ['tools', 22],
        ]);


        var options = {
          title: 'Release commits by project',
          is3D: true,
        };

        var chart = new
google.visualization.PieChart(document.getElementById('piechart'));
        chart.draw(data, options);
      }
    </script>
</head>
