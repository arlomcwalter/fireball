import java.io.FileNotFoundException;
import java.io.PrintWriter;
import java.io.UnsupportedEncodingException;
import java.util.List;

public class GenerateAst {
    public static void main(String[] args) throws FileNotFoundException, UnsupportedEncodingException {
        defineAst("src/main/java/minegame159/fireball", "Expr", List.of(
                "Binary   : Expr left, Token operator, Expr right",
                "Grouping : Expr expression",
                "Unary    : Token operator, Expr right"
        ));
    }

    private static void defineAst(String outputDir, String baseName, List<String> types) throws FileNotFoundException, UnsupportedEncodingException {
        PrintWriter writer = new PrintWriter(outputDir + "/" + baseName + ".java", "UTF-8");

        writer.println("package minegame159.fireball;");
        writer.println();
        writer.println("import java.util.List;");
        writer.println();
        writer.println("public abstract class " + baseName + " {");

        defineVisitor(writer, baseName, types);

        for (String type : types) {
            String className = type.split(":")[0].trim();
            String fields = type.split(":")[1].trim();

            writer.println();
            defineType(writer, baseName, className, fields);
        }

        writer.println();
        writer.println("    public abstract void accept(Visitor visitor);");

        writer.println("}");
        writer.close();
    }

    private static void defineVisitor(PrintWriter writer, String baseName, List<String> types) {
        writer.println("    interface Visitor {");

        for (String type : types) {
            String typeName = type.split(":")[0].trim();
            writer.println("        void visit" + typeName + baseName + "(" + typeName + " " + baseName.toLowerCase() + ");");
        }

        writer.println("    }");
    }

    private static void defineType(PrintWriter writer, String baseName, String className, String fieldList) {
        writer.println("    public static class " + className + " extends " + baseName + " {");

        // Fields
        String[] fields = fieldList.split(", ");

        for (String field : fields) {
            writer.println("        public final " + field + ";");
        }

        // Constructor
        writer.println();
        writer.println("        " + className + "(" + fieldList + ") {");

        for (String field : fields) {
            String name = field.split(" ")[1];
            writer.println("            this." + name + " = " + name + ";");
        }
        writer.println("        }");

        // Visitor pattern.
        writer.println();
        writer.println("        @Override");
        writer.println("        public void accept(Visitor visitor) {");
        writer.println("            visitor.visit" + className + baseName + "(this);");
        writer.println("        }");

        writer.println("    }");
    }
}
