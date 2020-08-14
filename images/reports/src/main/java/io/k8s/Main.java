package io.k8s;

import java.io.File;
import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.nio.file.StandardCopyOption;
import java.util.List;
import java.util.stream.Collectors;
import java.util.stream.Stream;
import net.masterthought.cucumber.Configuration;
import net.masterthought.cucumber.ReportBuilder;
import net.masterthought.cucumber.presentation.PresentationMode;
import net.masterthought.cucumber.sorting.SortingMethod;

public class Main {

  public static void main(String[] args) throws IOException, InterruptedException {
    String outputDirectory = System.getenv("OUTPUT_DIRECTORY");
    if (outputDirectory == null) {
      throw new RuntimeException("Environment variable OUTPUT_DIRECTORY is not optional");
    }

    File reportOutputDirectory = new File(outputDirectory);

    String inputFiles = System.getenv("INPUT_DIRECTORY");
    if (inputFiles == null) {
      throw new RuntimeException("Environment variable INPUT_DIRECTORY is not optional");
    }

    String build = System.getenv("BUILD");
    if (build == null) {
      throw new RuntimeException("Environment variable BUILD is not optional");
    }

    Stream<Path> walk = Files.walk(Paths.get(inputFiles));
    List<String> reports =
        walk.map(x -> x.toString())
            .filter(f -> f.endsWith("-report.json"))
            .collect(Collectors.toList());
    walk.close();

    String projectName = "Ingress Conformance";

    String release = System.getenv("RELEASE");

    Configuration configuration = new Configuration(reportOutputDirectory, projectName);
    configuration.setBuildNumber(build);
    configuration.setSortingMethod(SortingMethod.NATURAL);
    configuration.addPresentationModes(PresentationMode.EXPAND_ALL_STEPS);

    if (release != null) {
      configuration.addClassifications("Release", release);
    }

    String trendJSON =
        Path.of(outputDirectory, ReportBuilder.BASE_DIRECTORY, "trend.json").toString();
    configuration.setTrendsStatsFile(new File(trendJSON));

    ReportBuilder reportBuilder = new ReportBuilder(reports, configuration);
    reportBuilder.generateReports();

    Files.copy(
        Path.of(outputDirectory, ReportBuilder.BASE_DIRECTORY, "overview-features.html")
            .toFile()
            .toPath(),
        Path.of(outputDirectory, ReportBuilder.BASE_DIRECTORY, "index.html"), //
        StandardCopyOption.REPLACE_EXISTING //
        );

    // this should not that hard
    String baseDir = Path.of(outputDirectory, ReportBuilder.BASE_DIRECTORY).toString();
    ProcessBuilder builder = new ProcessBuilder();
    builder.command(
        "sh", "-c", "cp -R " + baseDir + "/* " + outputDirectory + " && rm -rf " + baseDir);
    Process process = builder.start();
    int exitCode = process.waitFor();
    assert exitCode == 0;
  }
}
