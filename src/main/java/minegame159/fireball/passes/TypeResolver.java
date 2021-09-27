package minegame159.fireball.passes;

import minegame159.fireball.Error;
import minegame159.fireball.Errors;
import minegame159.fireball.context.Context;
import minegame159.fireball.context.Field;
import minegame159.fireball.context.Function;
import minegame159.fireball.context.Struct;
import minegame159.fireball.parser.Expr;
import minegame159.fireball.parser.Parser;
import minegame159.fireball.parser.Stmt;
import minegame159.fireball.parser.Token;
import minegame159.fireball.types.StructType;
import minegame159.fireball.types.Type;

import java.util.*;

public class TypeResolver extends AstPass {
    private final List<Error> errors = new ArrayList<>();

    private final Context context;

    private final Stack<Map<String, Type>> scopes = new Stack<>();

    private TypeResolver(Context context) {
        this.context = context;
    }

    public static List<Error> resolve(Parser.Result result, Context context) {
        TypeResolver resolver = new TypeResolver(context);
        resolver.acceptS(result.stmts);
        return resolver.errors;
    }

    // Statements

    @Override
    public void visitExpressionStmt(Stmt.Expression stmt) {
        acceptE(stmt.expression);
    }

    @Override
    public void visitBlockStmt(Stmt.Block stmt) {
        scopes.push(new HashMap<>());
        acceptS(stmt.statements);
        scopes.pop();
    }

    @Override
    public void visitVariableStmt(Stmt.Variable stmt) {
        acceptE(stmt.initializer);

        Type type = stmt.getType(context);
        if (type == null) errors.add(Errors.unknownType(stmt.type, stmt.type));

        scopes.peek().put(stmt.name.lexeme(), type);
    }

    @Override
    public void visitIfStmt(Stmt.If stmt) {
        acceptE(stmt.condition);
        acceptS(stmt.thenBranch);
        acceptS(stmt.elseBranch);
    }

    @Override
    public void visitWhileStmt(Stmt.While stmt) {
        acceptE(stmt.condition);
        acceptS(stmt.body);
    }

    @Override
    public void visitForStmt(Stmt.For stmt) {
        acceptS(stmt.initializer);
        acceptE(stmt.condition);
        acceptE(stmt.increment);
        acceptS(stmt.body);
    }

    @Override
    public void visitFunctionStmt(Stmt.Function stmt) {
        acceptS(stmt.body);
    }

    @Override
    public void visitReturnStmt(Stmt.Return stmt) {
        acceptE(stmt.value);
    }

    @Override
    public void visitCBlockStmt(Stmt.CBlock stmt) {}

    @Override
    public void visitStructStmt(Stmt.Struct stmt) {}

    // Expressions

    @Override
    public void visitNullExpr(Expr.Null expr) {}

    @Override
    public void visitBoolExpr(Expr.Bool expr) {}

    @Override
    public void visitUnsignedIntExpr(Expr.UnsignedInt expr) {}

    @Override
    public void visitIntExpr(Expr.Int expr) {}

    @Override
    public void visitFloatExpr(Expr.Float expr) {}

    @Override
    public void visitStringExpr(Expr.String expr) {}

    @Override
    public void visitGroupingExpr(Expr.Grouping expr) {
        acceptE(expr.expression);
    }

    @Override
    public void visitBinaryExpr(Expr.Binary expr) {
        acceptE(expr.left);
        acceptE(expr.right);
    }

    @Override
    public void visitUnaryExpr(Expr.Unary expr) {
        acceptE(expr.right);
    }

    @Override
    public void visitLogicalExpr(Expr.Logical expr) {
        acceptE(expr.left);
        acceptE(expr.right);
    }

    @Override
    public void visitVariableExpr(Expr.Variable expr) {
        expr.type = resolveType(expr.name);
    }

    @Override
    public void visitAssignExpr(Expr.Assign expr) {
        expr.type = resolveType(expr.name);

        acceptE(expr.value);
    }

    @Override
    public void visitCallExpr(Expr.Call expr) {
        acceptE(expr.callee);
        acceptE(expr.arguments);
    }

    @Override
    public void visitGetExpr(Expr.Get expr) {
        acceptE(expr.object);

        expr.type = resolveFieldType(expr.object, expr.name);
    }

    @Override
    public void visitSetExpr(Expr.Set expr) {
        acceptE(expr.object);
        acceptE(expr.value);

        expr.type = resolveFieldType(expr.object, expr.name);
    }

    // Utils

    private Type resolveFieldType(Expr object, Token name) {
        if (!(object.getType() instanceof StructType)) errors.add(Errors.invalidFieldTarget(name));
        else {
            Struct struct = ((StructType) object.getType()).struct;
            Field field = struct.getField(name);

            if (field == null) errors.add(Errors.unknownField(struct.name(), name));
            else return field.type();
        }

        return null;
    }

    private Type resolveType(Token name) {
        Type type = getLocalType(name);

        if (type == null) {
            Function function = context.getFunction(name);
            if (function != null) type = function.returnType();
        }

        if (type == null) {
            errors.add(Errors.couldNotGetType(name));
        }

        return type;
    }

    private Type getLocalType(Token name) {
        for (int i = scopes.size() - 1; i >= 0; i--) {
            Type type = scopes.get(i).get(name.lexeme());
            if (type != null) return type;
        }

        return null;
    }
}
