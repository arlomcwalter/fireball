package minegame159.fireball.parser.prototypes;

import minegame159.fireball.parser.Stmt;
import minegame159.fireball.parser.Token;

import java.util.List;

public record ProtoStruct(Token name, List<ProtoParameter> fields, List<ProtoMethod> methods) {
    public ProtoStruct {
        for (ProtoMethod method : methods) method.owner = this;
    }

    public void accept(Stmt.Visitor visitor) {
        for (ProtoFunction method : methods) method.accept(visitor);
    }
}
